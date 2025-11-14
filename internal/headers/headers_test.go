package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test: Valid single header

func TestRequestLineParse(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("             Bearer: auth                    \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "auth", headers["bearer"])
	assert.Equal(t, 47, n)
	assert.False(t, done)

	// Test: valid done
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.True(t, done)

	// Test: valid two headers with existing values
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\n\r\n")
	otherData := []byte("Host: localhost:69420\r\n\r\n")
	n, done, err = headers.Parse(data)
	m, anotherDone, anotherErr := headers.Parse(otherData)
	require.NoError(t, err)
	require.NoError(t, anotherErr)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069, localhost:69420", headers["host"])
	assert.Equal(t, 23+23, n+m)
	assert.False(t, done)
	assert.False(t, anotherDone)

	// Test: valid two headers with different values
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\n\r\n")
	otherData = []byte("Bearer: auth\r\n\r\n")
	n, done, err = headers.Parse(data)
	m, anotherDone, anotherErr = headers.Parse(otherData)
	require.NoError(t, err)
	require.NoError(t, anotherErr)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "auth", headers["bearer"])
	assert.Equal(t, 23+14, n+m)
	assert.False(t, done)
	assert.False(t, anotherDone)

	// Test: invalid character in header key
	headers = NewHeaders()
	data = []byte("H)st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
