package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	str := string(data)

	crlfIndex := strings.Index(str, "\r\n")
	if crlfIndex == -1 {
		return 0, false, nil
	}

	// blank line => end of headers
	if crlfIndex == 0 {
		return 2, true, nil
	}

	line := str[:crlfIndex]

	colonIndex := strings.Index(line, ":")
	if colonIndex == -1 {
		return 0, false, errors.New("invalid header")
	}

	key := line[:colonIndex]

	// spaces before colon are invalid
	if strings.TrimSpace(key) != key {
		return 0, false, errors.New("invalid header spacing")
	}

	value := strings.TrimSpace(line[colonIndex+1:])

	h[key] = value

	return crlfIndex + 2, false, nil
}