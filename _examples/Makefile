.PHONY: all generate diff

all: generate diff

generate:
	cd golang-basics && make generate
	cd golang-imports && make generate

diff:
	git diff --color --ignore-all-space --ignore-blank-lines --exit-code .
