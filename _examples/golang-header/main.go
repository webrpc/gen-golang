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

func (rpc *ExampleServiceRPC) GetUser(ctx context.Context, req GetUserRequest) (*GetUserResponse, error) {
	if req.AuthToken == "" {
		return nil, ErrUnauthorized.WithCausef("missing auth token")
	}
	if req.UserID == 0 {
		return nil, ErrUserNotFound
	}
	return &GetUserResponse{
		User: &User{Id: req.UserID, Username: "alice"},
	}, nil
}

func (rpc *ExampleServiceRPC) CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
	if req.AuthToken == "" {
		return nil, ErrUnauthorized.WithCausef("missing auth token")
	}
	return &CreateUserResponse{
		User: &User{Id: 1, Username: req.Username},
	}, nil
}

func (rpc *ExampleServiceRPC) DeleteUser(ctx context.Context, req DeleteUserRequest) error {
	if req.AuthToken == "" {
		return ErrUnauthorized.WithCausef("missing auth token")
	}
	return nil
}
