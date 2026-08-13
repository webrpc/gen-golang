package main

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NewClients builds one working client per service from a shared address and
// HTTP client, so a caller talking to several services wires them in one call
// instead of a NewXClient(addr, http) line each.
func TestNewClients(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()
	go http.Serve(ln, newHandler())

	clients := NewClients("http://"+ln.Addr().String(), &http.Client{Timeout: 2 * time.Second})

	require.NotNil(t, clients.ExampleAPI)
	assert.NoError(t, clients.ExampleAPI.Ping(context.Background()))
}
