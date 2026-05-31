package request

import (
	"errors"
	"io"
	"regexp"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\r\n")

	requestLine, err := parseRequestLine(lines[0])
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: requestLine,
	}, nil
}

func parseRequestLine(line string) (RequestLine, error) {
	parts := strings.Split(line, " ")

	if len(parts) != 3 {
		return RequestLine{}, errors.New("invalid request line")
	}

	method := parts[0]
	target := parts[1]
	version := parts[2]

	methodRegex := regexp.MustCompile(`^[A-Z]+$`)
	if !methodRegex.MatchString(method) {
		return RequestLine{}, errors.New("invalid method")
	}

	if version != "HTTP/1.1" {
		return RequestLine{}, errors.New("invalid version")
	}

	return RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   "1.1",
	}, nil
}