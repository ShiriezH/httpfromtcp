package request

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLineParse(t *testing.T) {
	// Good GET Request line
	r, err := RequestFromReader(strings.NewReader(
		"GET / HTTP/1.1\r\nHost: localhost:42069\r\n\r\n",
	))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Good GET Request line with path
	r, err = RequestFromReader(strings.NewReader(
		"GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\n\r\n",
	))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Good POST Request with path
	r, err = RequestFromReader(strings.NewReader(
		"POST /coffee HTTP/1.1\r\nHost: localhost:42069\r\n\r\n",
	))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "POST", r.RequestLine.Method)

	// Invalid number of parts
	_, err = RequestFromReader(strings.NewReader(
		"/coffee HTTP/1.1\r\nHost: localhost:42069\r\n\r\n",
	))
	require.Error(t, err)

	// Invalid method
	_, err = RequestFromReader(strings.NewReader(
		"GeT /coffee HTTP/1.1\r\nHost: localhost:42069\r\n\r\n",
	))
	require.Error(t, err)

	// Invalid version
	_, err = RequestFromReader(strings.NewReader(
		"GET /coffee HTTP/2.0\r\nHost: localhost:42069\r\n\r\n",
	))
	require.Error(t, err)
}