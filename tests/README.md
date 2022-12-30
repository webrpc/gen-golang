# Interoperability tests

This folder implements the webrpc interoperability tests.

Since [webrpc](github.com/webrpc/webrpc) itself is written in Go, we have copied test/server test suite from the upstream repository.

Run `./update.sh` to update the test suite.

## Self-test

This unit test is ensuring that the client/server code generated from this repository can talk to each other.

1. Generate code
2. Test generated client/server code against each other

```bash
$ go test
```

## Test matrix against webrpc-test reference binary

These tests are ensuring that
- client code generated from this repository can talk to reference `webrpc-test@version` server
- server code generated from this repository responds to reference `webrpc-test@version` client tests

1. Generate code
2. Test generated client and server code against multiple versions of `webrpc-test` reference binaries via [test.sh](./test.sh) script

```bash
for webrpc in v0.10.0 v0.11.0; do
    ./test.sh $webrpc
done
```
