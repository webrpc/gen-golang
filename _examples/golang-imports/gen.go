package main

//go:generate go run github.com/webrpc/webrpc/cmd/webrpc-gen -schema=./proto/api.ridl -target=../../../gen-golang -out=./api.gen.go -pkg=main -server -client -legacyErrors=true -fmt=false

import (
	_ "github.com/webrpc/webrpc"
)
