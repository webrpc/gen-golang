package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

var client ExampleAPIClient

func TestMain(m *testing.M) {
	// Bind an ephemeral port first: the listener accepts connections before
	// http.Serve runs, so there is no startup race and nothing to wait for.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	go http.Serve(ln, newHandler())

	client = NewExampleAPIClient("http://"+ln.Addr().String(), &http.Client{
		Timeout: time.Duration(2 * time.Second),
	})

	code := m.Run()
	ln.Close()
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

// TestGetUser exercises the succinct request/response path end-to-end:
// the succinct client method and the server's succinctHandler.
func TestGetUser(t *testing.T) {
	resp, err := client.GetUser(context.Background(), GetUserRequest{Username: "alice"})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "alice", resp.User.Username)
	assert.Equal(t, uint32(30), resp.User.Age)
}

func TestListUsers(t *testing.T) {
	resp, err := client.ListUsers(context.Background(), ListUsersRequest{Page: 1})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Users, 1)
	assert.Equal(t, "pk", resp.Users[0].Username)
}

type registeredRoute struct {
	pattern string
	handler http.Handler
}

type recordingRouter struct {
	routes []registeredRoute
}

func (r *recordingRouter) Handle(pattern string, handler http.Handler) {
	r.routes = append(r.routes, registeredRoute{pattern: pattern, handler: handler})
}

type wrappedServer struct {
	WebRPCServer
}

func (s *wrappedServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.WebRPCServer.ServeHTTP(w, r)
}

func TestRegisterServer(t *testing.T) {
	server := NewExampleAPIServer(&ExampleRPC{})
	wantPatterns := []string{
		"/rpc/ExampleAPI/Ping",
		"/rpc/ExampleAPI/Status",
		"/rpc/ExampleAPI/GetUsers",
		"/rpc/ExampleAPI/GetUser",
		"/rpc/ExampleAPI/ListUsers",
	}

	t.Run("registers every exact route", func(t *testing.T) {
		router := &recordingRouter{}
		RegisterServer(router, server)

		patterns := make([]string, 0, len(router.routes))
		for _, route := range router.routes {
			patterns = append(patterns, route.pattern)
			assert.Same(t, server, route.handler)
		}
		assert.ElementsMatch(t, wantPatterns, patterns)
	})

	t.Run("dispatches through http ServeMux", func(t *testing.T) {
		router := http.NewServeMux()
		RegisterServer(router, server)

		req := httptest.NewRequest(http.MethodPost, "/rpc/ExampleAPI/Ping", strings.NewReader("{}"))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("preserves the Chi route pattern", func(t *testing.T) {
		var routePattern string
		router := chi.NewRouter()
		router.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
				routePattern = chi.RouteContext(r.Context()).RoutePattern()
			})
		})
		RegisterServer(router, server)

		req := httptest.NewRequest(http.MethodPost, "/rpc/ExampleAPI/GetUser", strings.NewReader(`{"username":"alice"}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "/rpc/ExampleAPI/GetUser", routePattern)
	})

	t.Run("registers a wrapped server", func(t *testing.T) {
		router := &recordingRouter{}
		wrapped := &wrappedServer{WebRPCServer: server}
		RegisterServer(router, wrapped)

		assert.Len(t, router.routes, len(wantPatterns))
		for _, route := range router.routes {
			assert.Same(t, wrapped, route.handler)
		}
	})
}
