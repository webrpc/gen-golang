# Interoperability tests

This folder implements the webrpc interoperability tests.

Since [webrpc](github.com/webrpc/webrpc) itself is written in Go, we have copied test/server test suite from the upstream repository.

Run `./update.sh` to update the test suite.

## Self-test

This unit test is ensuring that the client/server code generated from this repository can talk to each other.

```
$ go test
```

## Test matrix against webrpc-test reference binary

These tests are ensuring that
- client code generated from this repository can talk to reference `webrpc-test@version` server
- server code generated from this repository responds to reference `webrpc-test@version` client tests

Multiple versions of `webrpc-test` binaries can be tested from within [test.sh](./test.sh) script.

```
$ ./test.sh
```
