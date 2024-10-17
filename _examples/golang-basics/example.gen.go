// example v0.0.1 05b7a5c86b98738f4fe6ce9bb1fccd4af064847a
// --
// Code generated by webrpc-gen@v0.20.3-1-gf6584bc with ../../../gen-golang generator. DO NOT EDIT.
//
// webrpc-gen -schema=example.ridl -target=../../../gen-golang -pkg=main -server -client -legacyErrors -fixEmptyArrays -out=./example.gen.go
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
)

// WebRPC description and code-gen version
func WebRPCVersion() string {
	return "v1"
}

// Schema version of your RIDL schema
func WebRPCSchemaVersion() string {
	return "v0.0.1"
}

// Schema hash generated from your RIDL schema
func WebRPCSchemaHash() string {
	return "05b7a5c86b98738f4fe6ce9bb1fccd4af064847a"
}

//
// Common types
//

type Kind uint32

const (
	// user can only see number of transactions
	Kind_USER Kind = 0
	// admin permissions
	// can manage transactions
	// revert transactions
	// see analytics dashboard
	Kind_ADMIN Kind = 1
)

var Kind_name = map[uint32]string{
	0: "USER",
	1: "ADMIN",
}

var Kind_value = map[string]uint32{
	"USER":  0,
	"ADMIN": 1,
}

func (x Kind) String() string {
	return Kind_name[uint32(x)]
}

func (x Kind) MarshalText() ([]byte, error) {
	return []byte(Kind_name[uint32(x)]), nil
}

func (x *Kind) UnmarshalText(b []byte) error {
	*x = Kind(Kind_value[string(b)])
	return nil
}

func (x *Kind) Is(values ...Kind) bool {
	if x == nil {
		return false
	}
	for _, v := range values {
		if *x == v {
			return true
		}
	}
	return false
}

type Intent string

const (
	Intent_openSession     Intent = "openSession"
	Intent_closeSession    Intent = "closeSession"
	Intent_validateSession Intent = "validateSession"
)

func (x Intent) MarshalText() ([]byte, error) {
	return []byte(x), nil
}

func (x *Intent) UnmarshalText(b []byte) error {
	*x = Intent(string(b))
	return nil
}

func (x *Intent) Is(values ...Intent) bool {
	if x == nil {
		return false
	}
	for _, v := range values {
		if *x == v {
			return true
		}
	}
	return false
}

// Defines users within out wallet app
type User struct {
	ID   uint64    `json:"id" db:"id"`
	Uuid uuid.UUID `json:"uuid" db:"id"`
	// unique identifier of the user
	// must be unique !
	Username  string     `json:"USERNAME" db:"username"`
	Role      string     `json:"role" db:"-"`
	Nicknames []Nickname `json:"nicknames" db:"-"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
	Kind      Kind       `json:"kind"`
	Intent    Intent     `json:"intent"`
}

type Nickname struct {
	ID        uint64     `json:"ID" db:"id"`
	Nickname  string     `json:"nickname" db:"nickname"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
}

type SearchFilter struct {
	Q string `json:"q"`
}

type Version struct {
	WebrpcVersion string `json:"webrpcVersion"`
	SchemaVersion string `json:"schemaVersion"`
	SchemaHash    string `json:"schemaHash"`
}

type ComplexType struct {
	Meta              map[string]interface{}       `json:"meta"`
	MetaNestedExample map[string]map[string]uint32 `json:"metaNestedExample"`
	NamesList         []string                     `json:"namesList"`
	NumsList          []int64                      `json:"numsList"`
	DoubleArray       [][]string                   `json:"doubleArray"`
	ListOfMaps        []map[string]uint32          `json:"listOfMaps"`
	ListOfUsers       []*User                      `json:"listOfUsers"`
	MapOfUsers        map[string]*User             `json:"mapOfUsers"`
	User              *User                        `json:"user"`
}

var (
	methods = map[string]method{
		"/rpc/ExampleService/Ping": {
			Name:        "Ping",
			Service:     "ExampleService",
			Annotations: map[string]string{},
		},
		"/rpc/ExampleService/Status": {
			Name:        "Status",
			Service:     "ExampleService",
			Annotations: map[string]string{"internal": ""},
		},
		"/rpc/ExampleService/Version": {
			Name:        "Version",
			Service:     "ExampleService",
			Annotations: map[string]string{},
		},
		"/rpc/ExampleService/GetUser": {
			Name:        "GetUser",
			Service:     "ExampleService",
			Annotations: map[string]string{},
		},
		"/rpc/ExampleService/FindUser": {
			Name:        "FindUser",
			Service:     "ExampleService",
			Annotations: map[string]string{},
		},
		"/rpc/ExampleService/LogEvent": {
			Name:        "LogEvent",
			Service:     "ExampleService",
			Annotations: map[string]string{},
		},
	}
)

var WebRPCServices = map[string][]string{
	"ExampleService": {
		"Ping",
		"Status",
		"Version",
		"GetUser",
		"FindUser",
		"LogEvent",
	},
}

//
// Server types
//

type ExampleService interface {
	Ping(ctx context.Context) error
	// Status endpoint
	//
	// gives you current status of running application
	Status(ctx context.Context) (bool, error)
	Version(ctx context.Context) (*Version, error)
	GetUser(ctx context.Context, header map[string]string, userID uint64) (*User, error)
	FindUser(ctx context.Context, s *SearchFilter) (string, *User, error)
	LogEvent(ctx context.Context, event string) error
}

//
// Client types
//

type ExampleServiceClient interface {
	Ping(ctx context.Context) error
	// Status endpoint
	//
	// gives you current status of running application
	Status(ctx context.Context) (bool, error)
	Version(ctx context.Context) (*Version, error)
	GetUser(ctx context.Context, header map[string]string, userID uint64) (*User, error)
	FindUser(ctx context.Context, s *SearchFilter) (string, *User, error)
	LogEvent(ctx context.Context, event string) error
}

//
// Server
//

type WebRPCServer interface {
	http.Handler
}

type exampleServiceServer struct {
	ExampleService
	OnError   func(r *http.Request, rpcErr *WebRPCError)
	OnRequest func(w http.ResponseWriter, r *http.Request) error
}

func NewExampleServiceServer(svc ExampleService) *exampleServiceServer {
	return &exampleServiceServer{
		ExampleService: svc,
	}
}

func (s *exampleServiceServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		// In case of a panic, serve a HTTP 500 error and then panic.
		if rr := recover(); rr != nil {
			s.sendErrorJSON(w, r, ErrWebrpcServerPanic.WithCausef("%v", rr))
			panic(rr)
		}
	}()

	ctx := r.Context()
	ctx = context.WithValue(ctx, HTTPResponseWriterCtxKey, w)
	ctx = context.WithValue(ctx, HTTPRequestCtxKey, r)
	ctx = context.WithValue(ctx, ServiceNameCtxKey, "ExampleService")

	r = r.WithContext(ctx)

	var handler func(ctx context.Context, w http.ResponseWriter, r *http.Request)
	switch r.URL.Path {
	case "/rpc/ExampleService/Ping":
		handler = s.servePingJSON
	case "/rpc/ExampleService/Status":
		handler = s.serveStatusJSON
	case "/rpc/ExampleService/Version":
		handler = s.serveVersionJSON
	case "/rpc/ExampleService/GetUser":
		handler = s.serveGetUserJSON
	case "/rpc/ExampleService/FindUser":
		handler = s.serveFindUserJSON
	case "/rpc/ExampleService/LogEvent":
		handler = s.serveLogEventJSON
	default:
		err := ErrWebrpcBadRoute.WithCausef("no WebRPC method defined for path %v", r.URL.Path)
		s.sendErrorJSON(w, r, err)
		return
	}

	if r.Method != "POST" {
		w.Header().Add("Allow", "POST") // RFC 9110.
		err := ErrWebrpcBadMethod.WithCausef("unsupported method %v (only POST is allowed)", r.Method)
		s.sendErrorJSON(w, r, err)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if i := strings.Index(contentType, ";"); i >= 0 {
		contentType = contentType[:i]
	}
	contentType = strings.TrimSpace(strings.ToLower(contentType))

	switch contentType {
	case "application/json":
		if s.OnRequest != nil {
			if err := s.OnRequest(w, r); err != nil {
				rpcErr, ok := err.(WebRPCError)
				if !ok {
					rpcErr = ErrWebrpcEndpoint.WithCause(err)
				}
				s.sendErrorJSON(w, r, rpcErr)
				return
			}
		}

		handler(ctx, w, r)
	default:
		err := ErrWebrpcBadRequest.WithCausef("unsupported Content-Type %q (only application/json is allowed)", r.Header.Get("Content-Type"))
		s.sendErrorJSON(w, r, err)
	}
}

func (s *exampleServiceServer) servePingJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ctx = context.WithValue(ctx, MethodNameCtxKey, "Ping")

	// Call service method implementation.
	err := s.ExampleService.Ping(ctx)
	if err != nil {
		rpcErr, ok := err.(WebRPCError)
		if !ok {
			rpcErr = ErrWebrpcEndpoint.WithCause(err)
		}
		s.sendErrorJSON(w, r, rpcErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (s *exampleServiceServer) serveStatusJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ctx = context.WithValue(ctx, MethodNameCtxKey, "Status")

	// Call service method implementation.
	ret0, err := s.ExampleService.Status(ctx)
	if err != nil {
		rpcErr, ok := err.(WebRPCError)
		if !ok {
			rpcErr = ErrWebrpcEndpoint.WithCause(err)
		}
		s.sendErrorJSON(w, r, rpcErr)
		return
	}

	respPayload := struct {
		Ret0 bool `json:"status"`
	}{ret0}
	respBody, err := json.Marshal(initializeNilSlices(respPayload))
	if err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadResponse.WithCausef("failed to marshal json response: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func (s *exampleServiceServer) serveVersionJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ctx = context.WithValue(ctx, MethodNameCtxKey, "Version")

	// Call service method implementation.
	ret0, err := s.ExampleService.Version(ctx)
	if err != nil {
		rpcErr, ok := err.(WebRPCError)
		if !ok {
			rpcErr = ErrWebrpcEndpoint.WithCause(err)
		}
		s.sendErrorJSON(w, r, rpcErr)
		return
	}

	respPayload := struct {
		Ret0 *Version `json:"version"`
	}{ret0}
	respBody, err := json.Marshal(initializeNilSlices(respPayload))
	if err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadResponse.WithCausef("failed to marshal json response: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func (s *exampleServiceServer) serveGetUserJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ctx = context.WithValue(ctx, MethodNameCtxKey, "GetUser")

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadRequest.WithCausef("failed to read request data: %w", err))
		return
	}
	defer r.Body.Close()

	reqPayload := struct {
		Arg0 map[string]string `json:"header"`
		Arg1 uint64            `json:"userID"`
	}{}
	if err := json.Unmarshal(reqBody, &reqPayload); err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadRequest.WithCausef("failed to unmarshal request data: %w", err))
		return
	}

	// Call service method implementation.
	ret0, err := s.ExampleService.GetUser(ctx, reqPayload.Arg0, reqPayload.Arg1)
	if err != nil {
		rpcErr, ok := err.(WebRPCError)
		if !ok {
			rpcErr = ErrWebrpcEndpoint.WithCause(err)
		}
		s.sendErrorJSON(w, r, rpcErr)
		return
	}

	respPayload := struct {
		Ret0 *User `json:"user"`
	}{ret0}
	respBody, err := json.Marshal(initializeNilSlices(respPayload))
	if err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadResponse.WithCausef("failed to marshal json response: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func (s *exampleServiceServer) serveFindUserJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ctx = context.WithValue(ctx, MethodNameCtxKey, "FindUser")

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadRequest.WithCausef("failed to read request data: %w", err))
		return
	}
	defer r.Body.Close()

	reqPayload := struct {
		Arg0 *SearchFilter `json:"s"`
	}{}
	if err := json.Unmarshal(reqBody, &reqPayload); err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadRequest.WithCausef("failed to unmarshal request data: %w", err))
		return
	}

	// Call service method implementation.
	ret0, ret1, err := s.ExampleService.FindUser(ctx, reqPayload.Arg0)
	if err != nil {
		rpcErr, ok := err.(WebRPCError)
		if !ok {
			rpcErr = ErrWebrpcEndpoint.WithCause(err)
		}
		s.sendErrorJSON(w, r, rpcErr)
		return
	}

	respPayload := struct {
		Ret0 string `json:"name"`
		Ret1 *User  `json:"user"`
	}{ret0, ret1}
	respBody, err := json.Marshal(initializeNilSlices(respPayload))
	if err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadResponse.WithCausef("failed to marshal json response: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func (s *exampleServiceServer) serveLogEventJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ctx = context.WithValue(ctx, MethodNameCtxKey, "LogEvent")

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadRequest.WithCausef("failed to read request data: %w", err))
		return
	}
	defer r.Body.Close()

	reqPayload := struct {
		Arg0 string `json:"event"`
	}{}
	if err := json.Unmarshal(reqBody, &reqPayload); err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadRequest.WithCausef("failed to unmarshal request data: %w", err))
		return
	}

	// Call service method implementation.
	err = s.ExampleService.LogEvent(ctx, reqPayload.Arg0)
	if err != nil {
		rpcErr, ok := err.(WebRPCError)
		if !ok {
			rpcErr = ErrWebrpcEndpoint.WithCause(err)
		}
		s.sendErrorJSON(w, r, rpcErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (s *exampleServiceServer) sendErrorJSON(w http.ResponseWriter, r *http.Request, rpcErr WebRPCError) {
	if s.OnError != nil {
		s.OnError(r, &rpcErr)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rpcErr.HTTPStatus)

	respBody, _ := json.Marshal(rpcErr)
	w.Write(respBody)
}

func RespondWithError(w http.ResponseWriter, err error) {
	rpcErr, ok := err.(WebRPCError)
	if !ok {
		rpcErr = ErrWebrpcEndpoint.WithCause(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rpcErr.HTTPStatus)

	respBody, _ := json.Marshal(rpcErr)
	w.Write(respBody)
}

//
// Client
//

const ExampleServicePathPrefix = "/rpc/ExampleService/"

type exampleServiceClient struct {
	client HTTPClient
	urls   [6]string
}

func NewExampleServiceClient(addr string, client HTTPClient) ExampleServiceClient {
	prefix := urlBase(addr) + ExampleServicePathPrefix
	urls := [6]string{
		prefix + "Ping",
		prefix + "Status",
		prefix + "Version",
		prefix + "GetUser",
		prefix + "FindUser",
		prefix + "LogEvent",
	}
	return &exampleServiceClient{
		client: client,
		urls:   urls,
	}
}

func (c *exampleServiceClient) Ping(ctx context.Context) error {

	resp, err := doHTTPRequest(ctx, c.client, c.urls[0], nil, nil)
	if resp != nil {
		cerr := resp.Body.Close()
		if err == nil && cerr != nil {
			err = ErrWebrpcRequestFailed.WithCausef("failed to close response body: %w", cerr)
		}
	}

	return err
}

func (c *exampleServiceClient) Status(ctx context.Context) (bool, error) {
	out := struct {
		Ret0 bool `json:"status"`
	}{}

	resp, err := doHTTPRequest(ctx, c.client, c.urls[1], nil, &out)
	if resp != nil {
		cerr := resp.Body.Close()
		if err == nil && cerr != nil {
			err = ErrWebrpcRequestFailed.WithCausef("failed to close response body: %w", cerr)
		}
	}

	return out.Ret0, err
}

func (c *exampleServiceClient) Version(ctx context.Context) (*Version, error) {
	out := struct {
		Ret0 *Version `json:"version"`
	}{}

	resp, err := doHTTPRequest(ctx, c.client, c.urls[2], nil, &out)
	if resp != nil {
		cerr := resp.Body.Close()
		if err == nil && cerr != nil {
			err = ErrWebrpcRequestFailed.WithCausef("failed to close response body: %w", cerr)
		}
	}

	return out.Ret0, err
}

func (c *exampleServiceClient) GetUser(ctx context.Context, header map[string]string, userID uint64) (*User, error) {
	in := struct {
		Arg0 map[string]string `json:"header"`
		Arg1 uint64            `json:"userID"`
	}{header, userID}
	out := struct {
		Ret0 *User `json:"user"`
	}{}

	resp, err := doHTTPRequest(ctx, c.client, c.urls[3], in, &out)
	if resp != nil {
		cerr := resp.Body.Close()
		if err == nil && cerr != nil {
			err = ErrWebrpcRequestFailed.WithCausef("failed to close response body: %w", cerr)
		}
	}

	return out.Ret0, err
}

func (c *exampleServiceClient) FindUser(ctx context.Context, s *SearchFilter) (string, *User, error) {
	in := struct {
		Arg0 *SearchFilter `json:"s"`
	}{s}
	out := struct {
		Ret0 string `json:"name"`
		Ret1 *User  `json:"user"`
	}{}

	resp, err := doHTTPRequest(ctx, c.client, c.urls[4], in, &out)
	if resp != nil {
		cerr := resp.Body.Close()
		if err == nil && cerr != nil {
			err = ErrWebrpcRequestFailed.WithCausef("failed to close response body: %w", cerr)
		}
	}

	return out.Ret0, out.Ret1, err
}

func (c *exampleServiceClient) LogEvent(ctx context.Context, event string) error {
	in := struct {
		Arg0 string `json:"event"`
	}{event}

	resp, err := doHTTPRequest(ctx, c.client, c.urls[5], in, nil)
	if resp != nil {
		cerr := resp.Body.Close()
		if err == nil && cerr != nil {
			err = ErrWebrpcRequestFailed.WithCausef("failed to close response body: %w", cerr)
		}
	}

	return err
}

// HTTPClient is the interface used by generated clients to send HTTP requests.
// It is fulfilled by *(net/http).Client, which is sufficient for most users.
// Users can provide their own implementation for special retry policies.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// urlBase helps ensure that addr specifies a scheme. If it is unparsable
// as a URL, it returns addr unchanged.
func urlBase(addr string) string {
	// If the addr specifies a scheme, use it. If not, default to
	// http. If url.Parse fails on it, return it unchanged.
	url, err := url.Parse(addr)
	if err != nil {
		return addr
	}
	if url.Scheme == "" {
		url.Scheme = "http"
	}
	return url.String()
}

// newRequest makes an http.Request from a client, adding common headers.
func newRequest(ctx context.Context, url string, reqBody io.Reader, contentType string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", contentType)
	req.Header.Set("Content-Type", contentType)
	if headers, ok := HTTPRequestHeaders(ctx); ok {
		for k := range headers {
			for _, v := range headers[k] {
				req.Header.Add(k, v)
			}
		}
	}
	return req, nil
}

// doHTTPRequest is common code to make a request to the remote service.
func doHTTPRequest(ctx context.Context, client HTTPClient, url string, in, out interface{}) (*http.Response, error) {
	reqBody, err := json.Marshal(in)
	if err != nil {
		return nil, ErrWebrpcRequestFailed.WithCausef("failed to marshal JSON body: %w", err)
	}
	if err = ctx.Err(); err != nil {
		return nil, ErrWebrpcRequestFailed.WithCausef("aborted because context was done: %w", err)
	}

	req, err := newRequest(ctx, url, bytes.NewBuffer(reqBody), "application/json")
	if err != nil {
		return nil, ErrWebrpcRequestFailed.WithCausef("could not build request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, ErrWebrpcRequestFailed.WithCause(err)
	}

	if resp.StatusCode != 200 {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, ErrWebrpcBadResponse.WithCausef("failed to read server error response body: %w", err)
		}

		var rpcErr WebRPCError
		if err := json.Unmarshal(respBody, &rpcErr); err != nil {
			return nil, ErrWebrpcBadResponse.WithCausef("failed to unmarshal server error: %w", err)
		}
		if rpcErr.Cause != "" {
			rpcErr.cause = errors.New(rpcErr.Cause)
		}
		return nil, rpcErr
	}

	if out != nil {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, ErrWebrpcBadResponse.WithCausef("failed to read response body: %w", err)
		}

		err = json.Unmarshal(respBody, &out)
		if err != nil {
			return nil, ErrWebrpcBadResponse.WithCausef("failed to unmarshal JSON response body: %w", err)
		}
	}

	return resp, nil
}

func WithHTTPRequestHeaders(ctx context.Context, h http.Header) (context.Context, error) {
	if _, ok := h["Accept"]; ok {
		return nil, errors.New("provided header cannot set Accept")
	}
	if _, ok := h["Content-Type"]; ok {
		return nil, errors.New("provided header cannot set Content-Type")
	}

	copied := make(http.Header, len(h))
	for k, vv := range h {
		if vv == nil {
			copied[k] = nil
			continue
		}
		copied[k] = make([]string, len(vv))
		copy(copied[k], vv)
	}

	return context.WithValue(ctx, HTTPClientRequestHeadersCtxKey, copied), nil
}

func HTTPRequestHeaders(ctx context.Context) (http.Header, bool) {
	h, ok := ctx.Value(HTTPClientRequestHeadersCtxKey).(http.Header)
	return h, ok
}

//
// Helpers
//

type method struct {
	Name        string
	Service     string
	Annotations map[string]string
}

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "webrpc context value " + k.name
}

var (
	HTTPClientRequestHeadersCtxKey = &contextKey{"HTTPClientRequestHeaders"}
	HTTPResponseWriterCtxKey       = &contextKey{"HTTPResponseWriter"}

	HTTPRequestCtxKey = &contextKey{"HTTPRequest"}

	ServiceNameCtxKey = &contextKey{"ServiceName"}

	MethodNameCtxKey = &contextKey{"MethodName"}
)

func ServiceNameFromContext(ctx context.Context) string {
	service, _ := ctx.Value(ServiceNameCtxKey).(string)
	return service
}

func MethodNameFromContext(ctx context.Context) string {
	method, _ := ctx.Value(MethodNameCtxKey).(string)
	return method
}

func RequestFromContext(ctx context.Context) *http.Request {
	r, _ := ctx.Value(HTTPRequestCtxKey).(*http.Request)
	return r
}

func MethodCtx(ctx context.Context) (method, bool) {
	req := RequestFromContext(ctx)
	if req == nil {
		return method{}, false
	}

	m, ok := methods[req.URL.Path]
	if !ok {
		return method{}, false
	}

	return m, true
}

func ResponseWriterFromContext(ctx context.Context) http.ResponseWriter {
	w, _ := ctx.Value(HTTPResponseWriterCtxKey).(http.ResponseWriter)
	return w
}

// Copied from https://github.com/golang-cz/nilslice
func initializeNilSlices(obj interface{}) interface{} {
	v := reflect.ValueOf(obj)
	initializeNils(v)

	return obj
}

func initializeNils(v reflect.Value) {
	// Dereference pointer(s).
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	if v.Kind() == reflect.Slice {
		// Initialize a nil slice.
		if v.IsNil() && v.CanSet() {
			v.Set(reflect.MakeSlice(v.Type(), 0, 0))
			return
		}

		// Recursively iterate over slice items.
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i)
			initializeNils(item)
		}
	}

	// Recursively iterate over struct fields.
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			initializeNils(field)
		}
	}
}

//
// Errors
//

type WebRPCError struct {
	Name       string `json:"error"`
	Code       int    `json:"code"`
	Message    string `json:"msg"`
	Cause      string `json:"cause,omitempty"`
	HTTPStatus int    `json:"status"`
	cause      error
}

var _ error = WebRPCError{}

func (e WebRPCError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s %d: %s: %v", e.Name, e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("%s %d: %s", e.Name, e.Code, e.Message)
}

func (e WebRPCError) Is(target error) bool {
	if target == nil {
		return false
	}
	if rpcErr, ok := target.(WebRPCError); ok {
		return rpcErr.Code == e.Code
	}
	if legacyErr, ok := target.(legacyError); ok {
		return legacyErr.Code == e.Code
	}
	return errors.Is(e.cause, target)
}

func (e WebRPCError) Unwrap() error {
	return e.cause
}

func (e WebRPCError) WithCause(cause error) WebRPCError {
	err := e
	err.cause = cause
	err.Cause = cause.Error()
	return err
}

func (e WebRPCError) WithCausef(format string, args ...interface{}) WebRPCError {
	cause := fmt.Errorf(format, args...)
	err := e
	err.cause = cause
	err.Cause = cause.Error()
	return err
}

// Deprecated: Use .WithCause() method on WebRPCError.
func ErrorWithCause(rpcErr WebRPCError, cause error) WebRPCError {
	return rpcErr.WithCause(cause)
}

// Webrpc errors
var (
	ErrWebrpcEndpoint           = WebRPCError{Code: 0, Name: "WebrpcEndpoint", Message: "endpoint error", HTTPStatus: 400}
	ErrWebrpcRequestFailed      = WebRPCError{Code: -1, Name: "WebrpcRequestFailed", Message: "request failed", HTTPStatus: 400}
	ErrWebrpcBadRoute           = WebRPCError{Code: -2, Name: "WebrpcBadRoute", Message: "bad route", HTTPStatus: 404}
	ErrWebrpcBadMethod          = WebRPCError{Code: -3, Name: "WebrpcBadMethod", Message: "bad method", HTTPStatus: 405}
	ErrWebrpcBadRequest         = WebRPCError{Code: -4, Name: "WebrpcBadRequest", Message: "bad request", HTTPStatus: 400}
	ErrWebrpcBadResponse        = WebRPCError{Code: -5, Name: "WebrpcBadResponse", Message: "bad response", HTTPStatus: 500}
	ErrWebrpcServerPanic        = WebRPCError{Code: -6, Name: "WebrpcServerPanic", Message: "server panic", HTTPStatus: 500}
	ErrWebrpcInternalError      = WebRPCError{Code: -7, Name: "WebrpcInternalError", Message: "internal error", HTTPStatus: 500}
	ErrWebrpcClientDisconnected = WebRPCError{Code: -8, Name: "WebrpcClientDisconnected", Message: "client disconnected", HTTPStatus: 400}
	ErrWebrpcStreamLost         = WebRPCError{Code: -9, Name: "WebrpcStreamLost", Message: "stream lost", HTTPStatus: 400}
	ErrWebrpcStreamFinished     = WebRPCError{Code: -10, Name: "WebrpcStreamFinished", Message: "stream finished", HTTPStatus: 200}
)

// Schema errors
var (
	ErrMissingArgument = WebRPCError{Code: 500100, Name: "MissingArgument", Message: "missing argument", HTTPStatus: 400}
	ErrInvalidUsername = WebRPCError{Code: 500101, Name: "InvalidUsername", Message: "invalid username", HTTPStatus: 400}
	ErrMemoryFull      = WebRPCError{Code: 400100, Name: "MemoryFull", Message: "system memory is full", HTTPStatus: 400}
	ErrUnauthorized    = WebRPCError{Code: 400200, Name: "Unauthorized", Message: "unauthorized", HTTPStatus: 401}
	ErrUserNotFound    = WebRPCError{Code: 400300, Name: "UserNotFound", Message: "user not found", HTTPStatus: 400}
)

//
// Legacy errors
//

// Deprecated: Use fmt.Errorf() or WebRPCError.
func Errorf(err legacyError, format string, args ...interface{}) WebRPCError {
	return err.WebRPCError.WithCausef(format, args...)
}

// Deprecated: Use .WithCause() method on WebRPCError.
func WrapError(err legacyError, cause error, format string, args ...interface{}) WebRPCError {
	return err.WebRPCError.WithCausef("%v: %w", fmt.Errorf(format, args...), cause)
}

// Deprecated: Use fmt.Errorf() or WebRPCError.
func Failf(format string, args ...interface{}) WebRPCError {
	return Errorf(ErrFail, format, args...)
}

// Deprecated: Use .WithCause() method on WebRPCError.
func WrapFailf(cause error, format string, args ...interface{}) WebRPCError {
	return WrapError(ErrFail, cause, format, args...)
}

// Deprecated: Use fmt.Errorf() or WebRPCError.
func ErrorNotFound(format string, args ...interface{}) WebRPCError {
	return Errorf(ErrNotFound, format, args...)
}

// Deprecated: Use fmt.Errorf() or WebRPCError.
func ErrorInvalidArgument(argument string, validationMsg string) WebRPCError {
	return Errorf(ErrInvalidArgument, argument+" "+validationMsg)
}

// Deprecated: Use fmt.Errorf() or WebRPCError.
func ErrorRequiredArgument(argument string) WebRPCError {
	return ErrorInvalidArgument(argument, "is required")
}

// Deprecated: Use fmt.Errorf() or WebRPCError.
func ErrorInternal(format string, args ...interface{}) WebRPCError {
	return Errorf(ErrInternal, format, args...)
}

type legacyError struct{ WebRPCError }

// Legacy errors (webrpc v0.10.0 and earlier). Will be removed.
var (
	// Deprecated. Define errors in RIDL schema.
	ErrCanceled = legacyError{WebRPCError{Code: -10000, Name: "ErrCanceled", Message: "canceled", HTTPStatus: 408 /* RequestTimeout */}}
	// Deprecated. Define errors in RIDL schema.
	ErrUnknown = legacyError{WebRPCError{Code: -10001, Name: "ErrUnknown", Message: "unknown", HTTPStatus: 400 /* Bad Request */}}
	// Deprecated. Define errors in RIDL schema.
	ErrFail = legacyError{WebRPCError{Code: -10002, Name: "ErrFail", Message: "fail", HTTPStatus: 422 /* Unprocessable Entity */}}
	// Deprecated. Define errors in RIDL schema.
	ErrInvalidArgument = legacyError{WebRPCError{Code: -10003, Name: "ErrInvalidArgument", Message: "invalid argument", HTTPStatus: 400 /* BadRequest */}}
	// Deprecated. Define errors in RIDL schema.
	ErrDeadlineExceeded = legacyError{WebRPCError{Code: -10004, Name: "ErrDeadlineExceeded", Message: "deadline exceeded", HTTPStatus: 408 /* RequestTimeout */}}
	// Deprecated. Define errors in RIDL schema.
	ErrNotFound = legacyError{WebRPCError{Code: -10005, Name: "ErrNotFound", Message: "not found", HTTPStatus: 404 /* Not Found */}}
	// Deprecated. Define errors in RIDL schema.
	ErrBadRoute = legacyError{WebRPCError{Code: -10006, Name: "ErrBadRoute", Message: "bad route", HTTPStatus: 404 /* Not Found */}}
	// Deprecated. Define errors in RIDL schema.
	ErrAlreadyExists = legacyError{WebRPCError{Code: -10007, Name: "ErrAlreadyExists", Message: "already exists", HTTPStatus: 409 /* Conflict */}}
	// Deprecated. Define errors in RIDL schema.
	ErrPermissionDenied = legacyError{WebRPCError{Code: -10008, Name: "ErrPermissionDenied", Message: "permission denied", HTTPStatus: 403 /* Forbidden */}}
	// Deprecated. Define errors in RIDL schema.
	ErrUnauthenticated = legacyError{WebRPCError{Code: -10009, Name: "ErrUnauthenticated", Message: "unauthenticated", HTTPStatus: 401 /* Unauthorized */}}
	// Deprecated. Define errors in RIDL schema.
	ErrResourceExhausted = legacyError{WebRPCError{Code: -10010, Name: "ErrResourceExhausted", Message: "resource exhausted", HTTPStatus: 403 /* Forbidden */}}
	// Deprecated. Define errors in RIDL schema.
	ErrFailedPrecondition = legacyError{WebRPCError{Code: -10011, Name: "ErrFailedPrecondition", Message: "failed precondition", HTTPStatus: 412 /* Precondition Failed */}}
	// Deprecated. Define errors in RIDL schema.
	ErrAborted = legacyError{WebRPCError{Code: -10012, Name: "ErrAborted", Message: "aborted", HTTPStatus: 409 /* Conflict */}}
	// Deprecated. Define errors in RIDL schema.
	ErrOutOfRange = legacyError{WebRPCError{Code: -10013, Name: "ErrOutOfRange", Message: "out of range", HTTPStatus: 400 /* Bad Request */}}
	// Deprecated. Define errors in RIDL schema.
	ErrUnimplemented = legacyError{WebRPCError{Code: -10014, Name: "ErrUnimplemented", Message: "unimplemented", HTTPStatus: 501 /* Not Implemented */}}
	// Deprecated. Define errors in RIDL schema.
	ErrInternal = legacyError{WebRPCError{Code: -10015, Name: "ErrInternal", Message: "internal", HTTPStatus: 500 /* Internal Server Error */}}
	// Deprecated. Define errors in RIDL schema.
	ErrUnavailable = legacyError{WebRPCError{Code: -10016, Name: "ErrUnavailable", Message: "unavailable", HTTPStatus: 503 /* Service Unavailable */}}
	// Deprecated. Define errors in RIDL schema.
	ErrDataLoss = legacyError{WebRPCError{Code: -10017, Name: "ErrDataLoss", Message: "data loss", HTTPStatus: 500 /* Internal Server Error */}}
	// Deprecated. Define errors in RIDL schema.
	ErrNone = legacyError{WebRPCError{Code: -10018, Name: "ErrNone", Message: "", HTTPStatus: 200 /* OK */}}
)
