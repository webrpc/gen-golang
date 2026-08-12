package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// post calls a webrpc method and returns the raw response body. These tests
// assert on the JSON itself, because the wire format is what -fixEmptyArrays
// exists to control -- a typed client would hide the difference.
func post(t *testing.T, path, body string) string {
	t.Helper()

	r := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	newHandler().ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	return w.Body.String()
}

func TestEmptyArrays(t *testing.T) {
	t.Run("RequiredListsAreArraysWhenNil", func(t *testing.T) {
		// optional and label are absent rather than null, so a client can tell
		// "the server said nothing" apart from "the server said empty".
		assert.JSONEq(t, `{"report":{
			"required": [],
			"items":    [],
			"matrix":   [],
			"counts":   {},
			"buckets":  {},
			"byName":   {},
			"raw":      null,
			"count":    0
		}}`, post(t, "/rpc/EmptyArrays/GetReport", `{"id":"empty"}`))
	})

	t.Run("OptionalListIsArrayWhenSetToEmpty", func(t *testing.T) {
		// An explicitly empty optional list survives as [], and the nil lists
		// nested inside the item and the matrix are filled in too.
		assert.JSONEq(t, `{"report":{
			"required":  [],
			"optional":  [],
			"items":     [{"tags": []}],
			"matrix":    [[]],
			"counts":    {},
			"buckets":   {"a": []},
			"byName":    {"a": {"tags": []}},
			"optCounts": {},
			"raw":       null,
			"count":     0
		}}`, post(t, "/rpc/EmptyArrays/GetReport", `{"id":"explicit"}`))
	})

	t.Run("TopLevelReturnValueIsArrayWhenNil", func(t *testing.T) {
		// Multiple return values are marshaled through an anonymous struct,
		// which is the path a reflection-based fix cannot reach.
		assert.JSONEq(t, `{"reports":[],"total":0}`, post(t, "/rpc/EmptyArrays/ListReports", `{}`))
	})
}
