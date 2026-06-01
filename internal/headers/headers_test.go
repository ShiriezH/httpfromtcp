package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {

	// Valid single header
	headers := NewHeaders()
	data := []byte("HOST: localhost:42069\r\n\r\n")

	n, done, err := headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("HOST:      localhost:42069      \r\n\r\n")

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.False(t, done)

	// Valid 2 headers with existing headers
	headers = NewHeaders()
	headers["user-agent"] = "curl"

	data = []byte("HOST: localhost:42069\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, "curl", headers["user-agent"])
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.False(t, done)

	// Valid done
		// Valid duplicate header
	headers = NewHeaders()
	headers["set-person"] = "lane-loves-go"

	n, done, err = headers.Parse(
		[]byte("Set-Person: prime-loves-zig\r\n\r\n"),
	)

	require.NoError(t, err)
	assert.Equal(
		t,
		"lane-loves-go, prime-loves-zig",
		headers["set-person"],
	)
	assert.False(t, done)

	// Invalid spacing header
	headers = NewHeaders()

	n, done, err = headers.Parse(
		[]byte("       Host: localhost:42069\r\n\r\n"),
	)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Invalid character in header key
	headers = NewHeaders()

	n, done, err = headers.Parse(
		[]byte("H©st: localhost:42069\r\n\r\n"),
	)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}