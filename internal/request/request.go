package request

import (
	"errors"
	"io"
	"regexp"
	"strings"
)

const bufferSize = 8

const (
	stateInitialized = iota
	stateDone
)

type Request struct {
	RequestLine RequestLine
	state       int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0

	r := &Request{
		state: stateInitialized,
	}

	for r.state != stateDone {
		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:])

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		readToIndex += n

		parsed, err := r.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		if parsed > 0 {
			copy(buf, buf[parsed:readToIndex])
			readToIndex -= parsed
		}
	}

	if r.RequestLine.Method == "" {
		return nil, errors.New("failed to parse request line")
	}

	return r, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case stateInitialized:
		requestLine, consumed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if consumed == 0 {
			return 0, nil
		}

		r.RequestLine = requestLine
		r.state = stateDone

		return consumed, nil

	case stateDone:
		return 0, errors.New("trying to parse in done state")

	default:
		return 0, errors.New("unknown parser state")
	}
}

func parseRequestLine(data []byte) (RequestLine, int, error) {
	lineEnd := strings.Index(string(data), "\r\n")

	if lineEnd == -1 {
		return RequestLine{}, 0, nil
	}

	line := string(data[:lineEnd])
	parts := strings.Split(line, " ")

	if len(parts) != 3 {
		return RequestLine{}, 0, errors.New("invalid request line")
	}

	method := parts[0]
	target := parts[1]
	version := parts[2]

	methodRegex := regexp.MustCompile(`^[A-Z]+$`)
	if !methodRegex.MatchString(method) {
		return RequestLine{}, 0, errors.New("invalid method")
	}

	if version != "HTTP/1.1" {
		return RequestLine{}, 0, errors.New("invalid version")
	}

	return RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   "1.1",
	}, lineEnd + 2, nil
}