package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {

	// Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")

	n, done, err := headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("Host:      localhost:42069      \r\n\r\n")

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.False(t, done)

	// Valid 2 headers with existing headers
	headers = NewHeaders()
	headers["User-Agent"] = "curl"

	data = []byte("Host: localhost:42069\r\n\r\n")

	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.Equal(t, "curl", headers["User-Agent"])
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.False(t, done)

	// Valid done
	headers = NewHeaders()

	n, done, err = headers.Parse([]byte("\r\n"))

	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Invalid spacing header
	headers = NewHeaders()

	n, done, err = headers.Parse(
		[]byte("       Host: localhost:42069\r\n\r\n"),
	)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}