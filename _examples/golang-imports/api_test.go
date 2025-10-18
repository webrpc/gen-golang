package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var client ExampleAPIClient

func TestMain(m *testing.M) {
	go func() {
		if err := startServer(); err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(time.Millisecond * 500)

	client = NewExampleAPIClient("http://0.0.0.0:4242", &http.Client{
		Timeout: time.Duration(2 * time.Second),
	})

	code := m.Run()
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
