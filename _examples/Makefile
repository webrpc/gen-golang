.PHONY: all generate diff

all: generate diff

generate:
	webrpc-gen -schema=./golang-basics/example.ridl -target=../ -out=./golang-basics/example.gen.go -pkg=main -server -client
	webrpc-gen -schema=./golang-imports/proto/api.ridl -target=../ -out=./golang-imports/api.gen.go -pkg=main -server -client

diff:
	git diff --color --ignore-all-space --ignore-blank-lines --exit-code .