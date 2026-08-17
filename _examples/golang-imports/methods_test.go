package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Methods() exposes each RPC method bound to the server that handles it, so an
// application can register routes and attach per-method middleware itself.
func TestMethods(t *testing.T) {
	srv := NewExampleAPIServer(&ExampleRPC{})

	ms := srv.Methods()
	require.Len(t, ms, 5)

	byPath := make(map[string]Method, len(ms))
	for _, m := range ms {
		byPath[m.Path] = m
		assert.Equal(t, "ExampleAPI", m.Service())                    // promoted from *method
		assert.True(t, m.Handler == srv, "handler is the owning server")
		assert.True(t, strings.HasPrefix(m.Path, "/rpc/ExampleAPI/"))
	}

	// Annotations are reachable per method (Ping carries @auth in the schema).
	assert.True(t, byPath["/rpc/ExampleAPI/Ping"].HasAnnotation("auth"))
	assert.False(t, byPath["/rpc/ExampleAPI/Status"].HasAnnotation("auth"))
}

// Methods() is router-neutral: a stdlib *http.ServeMux is enough to mount the
// routes, and per-method middleware is applied by the application from
// annotations — no router-specific code in the generated output.
func TestMethodsMountRouterNeutral(t *testing.T) {
	srv := NewExampleAPIServer(&ExampleRPC{})

	var authed []string
	requireAuth := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authed = append(authed, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}

	mux := http.NewServeMux()
	for _, m := range srv.Methods() {
		h := m.Handler
		if m.HasAnnotation("auth") {
			h = requireAuth(h)
		}
		mux.Handle(m.Path, h)
	}

	post := func(path string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader("{}"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		return rec
	}

	assert.Equal(t, http.StatusOK, post("/rpc/ExampleAPI/Ping").Code)
	assert.Equal(t, http.StatusOK, post("/rpc/ExampleAPI/Status").Code)

	// The middleware ran only on the @auth method.
	assert.Equal(t, []string{"/rpc/ExampleAPI/Ping"}, authed)
}

// Server aggregates the Methods() of every non-nil service, so a whole API can be
// mounted from one struct literal; nil fields are skipped.
func TestServer(t *testing.T) {
	t.Run("aggregates a set service", func(t *testing.T) {
		ms := Server{ExampleAPI: &ExampleRPC{}}.Methods(nil)
		require.Len(t, ms, 5)
		for _, m := range ms {
			assert.Equal(t, "ExampleAPI", m.Service())
			assert.True(t, strings.HasPrefix(m.Path, "/rpc/ExampleAPI/"))
		}
	})

	t.Run("skips nil services", func(t *testing.T) {
		assert.Empty(t, Server{}.Methods(nil))
	})
}
