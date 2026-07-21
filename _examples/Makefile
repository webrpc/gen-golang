.PHONY: all generate diff test

all: generate diff test

generate:
	cd golang-basics && go generate -x ./...
	cd golang-imports && go generate -x ./...

diff:
	git diff --color --ignore-all-space --ignore-blank-lines --exit-code .

test:
	cd golang-basics && go test ./...
	cd golang-imports && go test ./...
