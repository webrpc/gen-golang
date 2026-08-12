package main

//go:generate go run github.com/webrpc/webrpc/cmd/webrpc-gen -schema=empty_arrays.ridl -target=../../../gen-golang -pkg=main -server -client -fixEmptyArrays=true -out=./empty_arrays.gen.go

import (
	"context"
	"log"
	"net/http"
)

func main() {
	if err := http.ListenAndServe(":4242", newHandler()); err != nil {
		log.Fatal(err)
	}
}

func newHandler() http.Handler {
	return NewEmptyArraysServer(&EmptyArraysRPC{})
}

type EmptyArraysRPC struct{}

// GetReport builds a report per the requested id, so a single succinct method
// covers both a fully unset report and one whose lists are set but empty.
func (s *EmptyArraysRPC) GetReport(ctx context.Context, req GetReportRequest) (*GetReportResponse, error) {
	if req.Id == "explicit" {
		return &GetReportResponse{Report: &Report{
			Optional:  []string{},
			Items:     []*Item{{}},
			Matrix:    [][]string{nil},
			Buckets:   map[string][]string{"a": nil},
			ByName:    map[string]*Item{"a": {}},
			OptCounts: map[string]uint32{},
		}}, nil
	}
	return &GetReportResponse{Report: &Report{}}, nil
}

// ListReports returns a nil list at the top level, alongside a second return
// value, so the response is built as an anonymous struct.
func (s *EmptyArraysRPC) ListReports(ctx context.Context) ([]*Report, uint32, error) {
	return nil, 0, nil
}
