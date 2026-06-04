package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
)

const port = 42069

const successHTML = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`

const badRequestHTML = `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`

const internalErrorHTML = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`

func main() {
	handler := func(
		w *response.Writer,
		req *request.Request,
	) {

		// NEW: video endpoint
		if req.RequestLine.RequestTarget == "/video" {
			videoData, err := os.ReadFile("assets/vim.mp4")
			if err != nil {
				return
			}

			respHeaders := response.GetDefaultHeaders(len(videoData))
			respHeaders.Set("Content-Type", "video/mp4")

			_ = w.WriteStatusLine(response.StatusOK)
			_ = w.WriteHeaders(respHeaders)
			_, _ = w.WriteBody(videoData)

			return
		}

		// httpbin proxy
		if strings.HasPrefix(
			req.RequestLine.RequestTarget,
			"/httpbin/",
		) {

			path := strings.TrimPrefix(
				req.RequestLine.RequestTarget,
				"/httpbin",
			)

			resp, err := http.Get(
				"https://httpbin.org" + path,
			)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			respHeaders := response.GetDefaultHeaders(0)

			respHeaders.Delete("Content-Length")
			respHeaders.Set("Transfer-Encoding", "chunked")
			respHeaders.Set(
				"Trailer",
				"X-Content-SHA256, X-Content-Length",
			)

			contentType := resp.Header.Get("Content-Type")
			if contentType == "" {
				contentType = "text/html"
			}

			respHeaders.Set("Content-Type", contentType)

			_ = w.WriteStatusLine(response.StatusOK)
			_ = w.WriteHeaders(respHeaders)

			buf := make([]byte, 1024)
			var fullBody []byte

			for {
				n, err := resp.Body.Read(buf)

				if n > 0 {
					chunk := buf[:n]
					fullBody = append(fullBody, chunk...)
					_, _ = w.WriteChunkedBody(chunk)
				}

				if err == io.EOF {
					break
				}

				if err != nil {
					return
				}
			}

			_, _ = w.WriteChunkedBodyDone()

			hash := sha256.Sum256(fullBody)

			trailers := headers.NewHeaders()
			trailers.Set(
				"X-Content-SHA256",
				hex.EncodeToString(hash[:]),
			)
			trailers.Set(
				"X-Content-Length",
				strconv.Itoa(len(fullBody)),
			)

			_ = w.WriteTrailers(trailers)

			return
		}

		var (
			status response.StatusCode
			body   string
		)

		switch req.RequestLine.RequestTarget {

		case "/yourproblem":
			status = response.StatusBadRequest
			body = badRequestHTML

		case "/myproblem":
			status = response.StatusInternalServerError
			body = internalErrorHTML

		default:
			status = response.StatusOK
			body = successHTML
		}

		respHeaders := response.GetDefaultHeaders(len(body))
		respHeaders.Set("Content-Type", "text/html")

		_ = w.WriteStatusLine(status)
		_ = w.WriteHeaders(respHeaders)
		_, _ = w.WriteBody([]byte(body))
	}

	srv, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	defer srv.Close()

	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	log.Println("Server gracefully stopped")
}