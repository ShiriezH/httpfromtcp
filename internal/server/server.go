package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync/atomic"

	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (h *HandlerError) Write(w io.Writer) error {
	err := response.WriteStatusLine(w, h.StatusCode)
	if err != nil {
		return err
	}

	body := []byte(h.Message)

	headers := response.GetDefaultHeaders(len(body))

	err = response.WriteHeaders(w, headers)
	if err != nil {
		return err
	}

	_, err = w.Write(body)
	return err
}

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: listener,
		handler:  handler,
	}

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		handlerErr := &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    "Bad Request\n",
		}
		_ = handlerErr.Write(conn)
		return
	}

	var body bytes.Buffer

	handlerErr := s.handler(&body, req)

	if handlerErr != nil {
		_ = handlerErr.Write(conn)
		return
	}

	bodyBytes := body.Bytes()

	err = response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		return
	}

	err = response.WriteHeaders(
		conn,
		response.GetDefaultHeaders(len(bodyBytes)),
	)
	if err != nil {
		return
	}

	_, _ = conn.Write(bodyBytes)
}