package main

//go:generate go run github.com/webrpc/webrpc/cmd/webrpc-gen -schema=example.ridl -target=../../../gen-golang -pkg=main -server -client -json=sonic -out=./example.gen.go

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

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

	webrpcHandler := NewExampleServer(&ExampleServiceRPC{})
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
	resp := ResponseWriterFromContext(ctx)
	serverVersions, err := VersionFromHeader(resp.Header())
	if err != nil {
		return nil, fmt.Errorf("parse server webrpc gen versions: %w", err)
	}

	req := RequestFromContext(ctx)
	clientVersions, err := VersionFromHeader(req.Header)
	if err != nil {
		return nil, fmt.Errorf("parse client webrpc gen versions: %w", err)
	}

	return &Version{
		WebrpcVersion: WebRPCVersion(),
		SchemaVersion: WebRPCSchemaVersion(),
		SchemaHash:    WebRPCSchemaHash(),
		ClientGenVersion: &GenVersions{
			WebrpcGenVersion: clientVersions.WebrpcGenVersion,
			TmplTarget:       clientVersions.CodeGenName,
			TmplVersion:      clientVersions.CodeGenVersion,
			SchemaVersion:    clientVersions.CodeGenVersion,
		},
		ServerGenVersion: &GenVersions{
			WebrpcGenVersion: serverVersions.WebrpcGenVersion,
			TmplTarget:       serverVersions.CodeGenName,
			TmplVersion:      serverVersions.CodeGenVersion,
			SchemaVersion:    serverVersions.CodeGenVersion,
		},
	}, nil
}

func (s *ExampleServiceRPC) GetUser(ctx context.Context, header map[string]string, userID uint64) (*User, error) {
	if userID == 911 {
		return nil, ErrUserNotFound.WithCausef("unknown user id %d", userID)
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
		return "", nil, ErrMissingArgument.WithCausef("s search filter required")
	}
	name := s.Q
	return s.Q, &User{
		ID:       123,
		Username: name,
	}, nil
}

func (rpc *ExampleServiceRPC) LogEvent(ctx context.Context, event string) error {
	return nil
}

func (rpc *ExampleServiceRPC) GetArticle(ctx context.Context, req GetArticleRequest) (*GetArticleResponse, error) {
	articleID := req.ArticleID
	if articleID == 0 {
		return nil, ErrMissingArgument.WithCausef("articleId is required")
	}
	if articleID == 404 {
		return nil, ErrUserNotFound.WithCausef("article not found")
	}

	content := "This is the content of the article."

	return &GetArticleResponse{
		Title:   fmt.Sprintf("Article %d", articleID),
		Content: &content,
	}, nil
}

func (rpc *ExampleServiceRPC) StreamNewArticles(ctx context.Context, stream StreamNewArticlesStreamWriter) error {
	for i := 0; i < 4; i++ {
		// content := fmt.Sprintf("This is the content of the article, %d", i)
		stream.Write(&GetArticleResponse{
			Title: fmt.Sprintf("Article %d", i),
			// Content: &content,
		})
		time.Sleep(1 * time.Second)
	}
	return nil
}
