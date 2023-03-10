package main

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	client ExampleService
)

// func TestMain()

func init() {
	go func() {
		startServer()
	}()

	client = NewExampleServiceClient("http://0.0.0.0:4242", &http.Client{
		Timeout: time.Duration(2 * time.Second),
	})
	time.Sleep(time.Millisecond * 500)

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

func TestGetUser(t *testing.T) {
	{
		arg1 := map[string]string{"a": "1"}
		user, err := client.GetUser(context.Background(), arg1, 12)
		assert.Equal(t, &User{ID: 12, Username: "hihi"}, user)
		assert.NoError(t, err)
	}

	{
		// Error case, expecting to receive an error
		user, err := client.GetUser(context.Background(), nil, 911)
		assert.True(t, errors.Is(err, ErrUserNotFound))
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrUserNotFound)
		assert.Contains(t, err.Error(), "not found")

		rpcErr, ok := err.(WebRPCError)
		assert.True(t, ok)
		assert.Contains(t, rpcErr.Unwrap().Error(), "911")
	}

	{
		name, user, err := client.FindUser(context.Background(), &SearchFilter{Q: "joe"})
		assert.Equal(t, "joe", name)
		assert.Equal(t, &User{ID: 123, Username: "joe"}, user)
		assert.NoError(t, err)
	}
}

func TestLegacyErrors(t *testing.T) {
	{
		_, err := client.GetUser(context.Background(), nil, 0)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)

		_, err = client.GetUser(context.Background(), nil, 1000)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrUnavailable)
	}
}
