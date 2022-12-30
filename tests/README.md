# Interoperability tests

This folder implements webrpc interoperability tests.

Since [webrpc](github.com/webrpc/webrpc) itself is written in Go and has interoperability tests, we have simply copied the test cases from the upstream repository via `./update.sh` script.

We compile the test cases against a fresh client/server code generated from this repository. To regenerate client/server code, run:

```bash
go generate ./...
```

To update the test cases, edit the `VERSION` variable in `./update.sh` script and run it:

```bash
./update.sh
```

## Self-test

The basic interoperability unit test is ensuring that the client and server code generated from this repository can talk to each other.

1. Generate code
2. Test generated client/server code against each other

```bash
go generate ./...
go test ./...
```

## Test matrix against webrpc-test reference binary

These tests are ensuring that
- client code generated from this repository can talk to reference `webrpc-test@version` server
- server code generated from this repository responds to reference `webrpc-test@version` client tests

1. Generate code
2. Test generated client and server code against multiple versions of `webrpc-test` reference binaries via [test.sh](./test.sh) script

```bash
go generate ./...

for webrpcVersion in v0.10.0; do
    ./download.sh $webrpcVersion bin/$webrpcVersion
    PATH="bin/$webrpcVersion:$PATH" ./test.sh
done
```
