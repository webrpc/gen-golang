.PHONY: all generate diff

all: generate diff

generate:
	cd golang-basics && go generate -x ./...
	cd golang-imports && go generate -x ./...

diff:
	git diff --color --ignore-all-space --ignore-blank-lines --exit-code .
