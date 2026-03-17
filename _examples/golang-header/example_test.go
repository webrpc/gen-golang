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
	resp, err := client.GetUser(context.Background(), GetUserRequest{
		AuthToken: "my-secret-token",
		UserID:    42,
	})
	require.NoError(t, err)
	assert.Equal(t, uint64(42), resp.User.Id)
	assert.Equal(t, "alice", resp.User.Username)
}

func TestGetUser_MissingToken(t *testing.T) {
	_, err := client.GetUser(context.Background(), GetUserRequest{
		AuthToken: "",
		UserID:    42,
	})
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrUnauthorized)
}

func TestCreateUser_MultipleHeaders(t *testing.T) {
	resp, err := client.CreateUser(context.Background(), CreateUserRequest{
		AuthToken: "my-secret-token",
		Role:      "admin",
		Username:  "bob",
	})
	require.NoError(t, err)
	assert.Equal(t, "bob", resp.User.Username)
}

func TestDeleteUser_AllHeaders(t *testing.T) {
	err := client.DeleteUser(context.Background(), DeleteUserRequest{
		AuthToken: "my-secret-token",
		UserID:    99,
	})
	assert.NoError(t, err)
}

func TestDeleteUser_Unauthorized(t *testing.T) {
	err := client.DeleteUser(context.Background(), DeleteUserRequest{
		AuthToken: "",
		UserID:    99,
	})
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrUnauthorized)
}
