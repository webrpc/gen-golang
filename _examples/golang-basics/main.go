//go:generate webrpc-gen -schema=example.ridl -target=../../../gen-golang -pkg=main -server -client -out=./example.gen.go -fmt=false
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	err := startServer()
	if err != nil {
		log.Fatal(err)
	}
}

func startServer() error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("."))
	})

	webrpcHandler := NewExampleServiceServer(&ExampleServiceRPC{})
	r.Handle("/*", webrpcHandler)

	log.Printf("Listening on :4242")
	return http.ListenAndServe(":4242", r)
}

type ExampleServiceRPC struct {
}

func (rpc *ExampleServiceRPC) Ping(ctx context.Context) error {
	return nil
}

func (rpc *ExampleServiceRPC) Status(ctx context.Context) (bool, error) {
	return true, nil
}

func (rpc *ExampleServiceRPC) Version(ctx context.Context) (*Version, error) {
	return &Version{
		WebrpcVersion: WebRPCVersion(),
		SchemaVersion: WebRPCSchemaVersion(),
		SchemaHash:    WebRPCSchemaHash(),
	}, nil
}

func (s *ExampleServiceRPC) GetUser(ctx context.Context, header map[string]string, userID uint64) (*User, error) {
	if userID == 911 {
		return nil, ErrorWithCause(ErrUserNotFound, fmt.Errorf("unknown user id %d", userID))
	}
	if userID == 31337 {
		return nil, ErrUnauthorized
	}
	if userID == 666 {
		panic("oh no")
	}

	return &User{
		ID:       userID,
		Username: "hihi",
	}, nil
}

func (rpc *ExampleServiceRPC) FindUser(ctx context.Context, s *SearchFilter) (string, *User, error) {
	if s == nil {
		return "", nil, ErrorWithCause(ErrMissingArgument, fmt.Errorf("s search filter required"))
	}
	name := s.Q
	return s.Q, &User{
		ID:       123,
		Username: name,
	}, nil
}
