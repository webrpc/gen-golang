package main

//go:generate go run github.com/webrpc/webrpc/cmd/webrpc-gen -schema=example.ridl -target=../../../gen-golang -pkg=main -server -client -out=./example.gen.go

import (
	"context"
	"log"
	"net/http"
)

func main() {
	err := startServer()
	if err != nil {
		log.Fatal(err)
	}
}

func startServer() error {
	webrpcHandler := NewExampleServiceServer(&ExampleServiceRPC{})

	log.Printf("Listening on :4243")
	return http.ListenAndServe(":4243", webrpcHandler)
}

type ExampleServiceRPC struct{}

func (rpc *ExampleServiceRPC) Ping(ctx context.Context) error {
	return nil
}

func (rpc *ExampleServiceRPC) GetUser(ctx context.Context, authToken string, userID uint64) (*User, error) {
	if authToken == "" {
		return nil, ErrUnauthorized.WithCausef("missing auth token")
	}
	if userID == 0 {
		return nil, ErrUserNotFound
	}
	return &User{
		Id:       userID,
		Username: "alice",
	}, nil
}

func (rpc *ExampleServiceRPC) CreateUser(ctx context.Context, authToken string, role string, username string) (*User, error) {
	if authToken == "" {
		return nil, ErrUnauthorized.WithCausef("missing auth token")
	}
	return &User{
		Id:       1,
		Username: username,
	}, nil
}

func (rpc *ExampleServiceRPC) DeleteUser(ctx context.Context, authToken string, userID uint64) error {
	if authToken == "" {
		return ErrUnauthorized.WithCausef("missing auth token")
	}
	return nil
}
