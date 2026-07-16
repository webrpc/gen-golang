.PHONY: all generate diff test

all: generate diff test

generate:
	cd golang-basics && go generate -x ./...
	cd golang-imports && go generate -x ./...

diff:
	git diff --color --ignore-all-space --ignore-blank-lines --exit-code .

# golang-basics is generate+diff only: its generated code maps a type to
# github.com/0xsequence/go-sequence/lib/prototyp, which the module does not
# require, so it does not build. golang-imports covers the runtime path.
test:
	cd golang-imports && go test ./...
