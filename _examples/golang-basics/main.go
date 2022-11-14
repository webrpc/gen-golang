//go:generate webrpc-gen -schema=example.ridl -target=../../../gen-golang -pkg=main -server -client -out=./example.gen.go
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

func (rpc *ExampleServiceRPC) GetUser(ctx context.Context, header map[string]string, userID uint64) (uint32, *User, error) {
	if userID == 911 {
		return 0, nil, Errorf(ErrNotFound, "unknown userID %d", 911)
	}

	return 200, &User{
		ID:       userID,
		Username: "hihi",
	}, nil
}

func (rpc *ExampleServiceRPC) FindUser(ctx context.Context, s *SearchFilter) (string, *User, error) {
	if s == nil {
		return "", nil, Errorf(ErrInvalidArgument, "s search filter required")
	}
	name := s.Q
	return s.Q, &User{
		ID:       123,
		Username: name,
	}, nil
}
