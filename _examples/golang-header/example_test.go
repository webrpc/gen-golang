package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	client ExampleServiceClient
)

func init() {
	go func() {
		startServer()
	}()

	client = NewExampleServiceClient("http://0.0.0.0:4243", &http.Client{
		Timeout: 10 * time.Second,
	})
	time.Sleep(time.Millisecond * 500)
}

func TestPing(t *testing.T) {
	err := client.Ping(context.Background())
	assert.NoError(t, err)
}

func TestGetUser_WithHeader(t *testing.T) {
	// authToken is sent as HTTP header, userID in JSON body
	user, err := client.GetUser(context.Background(), "my-secret-token", 42)
	require.NoError(t, err)
	assert.Equal(t, uint64(42), user.Id)
	assert.Equal(t, "alice", user.Username)
}

func TestGetUser_MissingToken(t *testing.T) {
	// Empty auth token should return unauthorized
	_, err := client.GetUser(context.Background(), "", 42)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrUnauthorized)
}

func TestCreateUser_MultipleHeaders(t *testing.T) {
	// authToken and role are sent as HTTP headers, username in JSON body
	user, err := client.CreateUser(context.Background(), "my-secret-token", "admin", "bob")
	require.NoError(t, err)
	assert.Equal(t, "bob", user.Username)
}

func TestDeleteUser_AllHeaders(t *testing.T) {
	// Both authToken and userID sent as HTTP headers, no JSON body
	err := client.DeleteUser(context.Background(), "my-secret-token", 99)
	assert.NoError(t, err)
}

func TestDeleteUser_Unauthorized(t *testing.T) {
	err := client.DeleteUser(context.Background(), "", 99)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrUnauthorized)
}
