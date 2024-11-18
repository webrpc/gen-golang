package main

//go:generate go run github.com/webrpc/webrpc/cmd/webrpc-gen -schema=./proto/api.ridl -target=../../../gen-golang -out=./api.gen.go -pkg=main -server -client -fmt=false

import (
	"context"
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

	webrpcHandler := NewExampleAPIServer(&ExampleRPC{})
	r.Handle("/*", webrpcHandler)

	return http.ListenAndServe(":4242", r)
}

type ExampleRPC struct {
}

func (s *ExampleRPC) Ping(ctx context.Context) error {
	return nil
}

func (s *ExampleRPC) Status(ctx context.Context) (bool, error) {
	return true, nil
}

func (s *ExampleRPC) GetUsers(ctx context.Context) ([]*User, Location, error) {
	loc := Location_TORONTO
	return []*User{
		{Username: "pk", Age: 99},
	}, loc, nil
}
