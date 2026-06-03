package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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

		headers := response.GetDefaultHeaders(len(body))
		headers.Set("Content-Type", "text/html")

		_ = w.WriteStatusLine(status)
		_ = w.WriteHeaders(headers)
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