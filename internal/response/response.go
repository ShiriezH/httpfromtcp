package response

import (
	"fmt"
	"io"
	"strconv"

	"httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	reasonPhrase := ""

	switch statusCode {
	case StatusOK:
		reasonPhrase = "OK"
	case StatusBadRequest:
		reasonPhrase = "Bad Request"
	case StatusInternalServerError:
		reasonPhrase = "Internal Server Error"
	}

	_, err := fmt.Fprintf(
		w.writer,
		"HTTP/1.1 %d %s\r\n",
		statusCode,
		reasonPhrase,
	)

	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	for key, value := range headers {
		_, err := fmt.Fprintf(
			w.writer,
			"%s: %s\r\n",
			key,
			value,
		)
		if err != nil {
			return err
		}
	}

	_, err := fmt.Fprint(w.writer, "\r\n")
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	_, err := fmt.Fprintf(
		w.writer,
		"%x\r\n",
		len(p),
	)
	if err != nil {
		return 0, err
	}

	n, err := w.writer.Write(p)
	if err != nil {
		return n, err
	}

	_, err = fmt.Fprint(w.writer, "\r\n")
	if err != nil {
		return n, err
	}

	return n, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	_, err := fmt.Fprint(w.writer, "0\r\n")
	if err != nil {
		return 0, err
	}

	return 3, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	for key, value := range h {
		_, err := fmt.Fprintf(
			w.writer,
			"%s: %s\r\n",
			key,
			value,
		)
		if err != nil {
			return err
		}
	}

	_, err := fmt.Fprint(w.writer, "\r\n")
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()

	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}