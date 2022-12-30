#!/bin/bash
set -e

OS=$(uname -o | tr A-Z a-z)
ARCH=$(uname -m)
PORT=9889

VERSION=v0.10.0

# Update test suite from https://github.com/webrpc/webrpc.
curl -o schema/test.ridl -LJO https://raw.githubusercontent.com/webrpc/webrpc/$VERSION/tests/schema/test.ridl
curl -o client/client.go -LJO https://raw.githubusercontent.com/webrpc/webrpc/$VERSION/tests/client/client.go
curl -o server/server.go -LJO https://raw.githubusercontent.com/webrpc/webrpc/$VERSION/tests/server/server.go

WEBRPC_GEN="bin/webrpc-gen@$VERSION.$OS-$ARCH"
WEBRPC_TEST="bin/webrpc-test@$VERSION.$OS-$ARCH"

# Download webrpc binaries if not available locally
WEBRPC_GEN_URL="https://github.com/webrpc/webrpc/releases/download/$VERSION/webrpc-gen.$OS-$ARCH"
WEBRPC_TEST_URL="https://github.com/webrpc/webrpc/releases/download/$VERSION/webrpc-test.$OS-$ARCH"
[[ ! -f $WEBRPC_GEN ]] && curl -o $WEBRPC_GEN -LJO "$WEBRPC_GEN_URL" && chmod +x $WEBRPC_GEN
[[ ! -f $WEBRPC_TEST ]] && curl -o $WEBRPC_TEST -LJO "$WEBRPC_TEST_URL" && chmod +x $WEBRPC_TEST

$WEBRPC_GEN -schema=./schema/test.ridl -target=golang -pkg=server -server -out=./server/server.gen.go
$WEBRPC_GEN -schema=./schema/test.ridl -target=golang -pkg=client -client -out=./client/client.gen.go
