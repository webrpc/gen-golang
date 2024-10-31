package main

//go:generate go run github.com/webrpc/webrpc/cmd/webrpc-gen -schema=example.ridl -target=../../../gen-golang -pkg=main -server -client -legacyErrors -fixEmptyArrays -out=./example.gen.go

import (
	_ "github.com/webrpc/webrpc"
)
