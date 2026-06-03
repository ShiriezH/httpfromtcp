package headers

import (
	"errors"
	"regexp"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Get(key string) string {
	return h[strings.ToLower(key)]
}

func (h Headers) Set(key, value string) {
	h[strings.ToLower(key)] = value
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	str := string(data)

	crlfIndex := strings.Index(str, "\r\n")
	if crlfIndex == -1 {
		return 0, false, nil
	}

	if crlfIndex == 0 {
		return 2, true, nil
	}

	line := str[:crlfIndex]

	colonIndex := strings.Index(line, ":")
	if colonIndex == -1 {
		return 0, false, errors.New("invalid header")
	}

	key := line[:colonIndex]

	if strings.TrimSpace(key) != key {
		return 0, false, errors.New("invalid header spacing")
	}

	validKey := regexp.MustCompile(`^[A-Za-z0-9!#$%&'*+\-.^_` + "`" + `|~]+$`)
	if !validKey.MatchString(key) {
		return 0, false, errors.New("invalid header key")
	}

	value := strings.TrimSpace(line[colonIndex+1:])

	key = strings.ToLower(key)

	if existingValue, exists := h[key]; exists {
		h[key] = existingValue + ", " + value
	} else {
		h[key] = value
	}

	return crlfIndex + 2, false, nil
}