package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var client ExampleAPIClient

func TestMain(m *testing.M) {
	// Bind an ephemeral port first: the listener accepts connections before
	// http.Serve runs, so there is no startup race and nothing to wait for.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	go http.Serve(ln, newHandler())

	client = NewExampleAPIClient("http://"+ln.Addr().String(), &http.Client{
		Timeout: time.Duration(2 * time.Second),
	})

	code := m.Run()
	ln.Close()
	os.Exit(code)
}

func TestPing(t *testing.T) {
	err := client.Ping(context.Background())
	assert.NoError(t, err)
}

func TestStatus(t *testing.T) {
	resp, err := client.Status(context.Background())
	assert.Equal(t, true, resp)
	assert.NoError(t, err)
}

// TestGetUser exercises the succinct request/response path end-to-end:
// the succinct client method and the server's succinctHandler.
func TestGetUser(t *testing.T) {
	resp, err := client.GetUser(context.Background(), GetUserRequest{Username: "alice"})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "alice", resp.User.Username)
	assert.Equal(t, uint32(30), resp.User.Age)
}

func TestListUsers(t *testing.T) {
	resp, err := client.ListUsers(context.Background(), ListUsersRequest{Page: 1})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Users, 1)
	assert.Equal(t, "pk", resp.Users[0].Username)
}
