#!/bin/bash
set -e

### Run interoperability tests against specific version of webrpc-test.

VERSION="${1}"

[[ -z "$VERSION" ]] && { echo "Usage: $0 <webrpc-version>"; exit 1; }

OS=$(uname -o | tr A-Z a-z)
ARCH=$(uname -m)
PORT=9889

WEBRPC_GEN="bin/webrpc-gen@$VERSION.$OS-$ARCH"
WEBRPC_TEST="bin/webrpc-test@$VERSION.$OS-$ARCH"

# Download webrpc binaries if not available locally
WEBRPC_GEN_URL="https://github.com/webrpc/webrpc/releases/download/$VERSION/webrpc-gen.$OS-$ARCH"
WEBRPC_TEST_URL="https://github.com/webrpc/webrpc/releases/download/$VERSION/webrpc-test.$OS-$ARCH"
[[ ! -f $WEBRPC_GEN ]] && curl -o $WEBRPC_GEN -fLJO "$WEBRPC_GEN_URL" && chmod +x $WEBRPC_GEN
[[ ! -f $WEBRPC_TEST ]] && curl -o $WEBRPC_TEST -fLJO "$WEBRPC_TEST_URL" && chmod +x $WEBRPC_TEST

echo "###"
echo "### webrpc@$VERSION reference server <=> generated client"
echo "###"
echo

# Run reference webrpc-test@VERSION server and wait for it to be ready
echo "Running reference webrpc-test@$VERSION server at 0.0.0.0:$PORT"
$WEBRPC_TEST -server -port=$PORT -timeout=2s &

# Wait until http://localhost:$PORT is available.
until nc -z localhost $PORT; do sleep 0.1; done

# Run generated client tests
go test -v -server=false -client=true -url=http://localhost:$PORT

wait
echo

echo "###"
echo "### generated server <=> webrpc@$VERSION reference client"
echo "###"
echo

# Run generated server
go test -v -server=true -client=false -httptest.serve=0.0.0.0:$PORT -serverTimeout=2s &
until nc -z localhost $PORT; do sleep 0.1; done

# Run reference webrpc-test@VERSION client and wait for it to be ready
echo "Running reference webrpc-test@$VERSION client tests against http://localhost:$PORT"
$WEBRPC_TEST -client -url=http://localhost:$PORT

wait
echo
