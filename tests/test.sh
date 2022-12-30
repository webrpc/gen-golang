#!/bin/bash
set -e

OS=$(uname -o | tr A-Z a-z)
ARCH=$(uname -m)
PORT=9889

### Run interoperability tests against specific versions of webrpc.
for VERSION in v0.9.0 v0.9.1; do
    WEBRPC_GEN="bin/webrpc-gen@$VERSION.$OS-$ARCH"
    WEBRPC_TEST="bin/webrpc-test@$VERSION.$OS-$ARCH"

    # Download webrpc binaries if not available locally
    WEBRPC_GEN_URL="https://github.com/webrpc/webrpc/releases/download/$VERSION/webrpc-gen.$OS-$ARCH"
    WEBRPC_TEST_URL="https://github.com/webrpc/webrpc/releases/download/$VERSION/webrpc-test.$OS-$ARCH"
    [[ ! -f $WEBRPC_GEN ]] && curl -o $WEBRPC_GEN -LJO "$WEBRPC_GEN_URL" && chmod +x $WEBRPC_GEN
    [[ ! -f $WEBRPC_TEST ]] && curl -o $WEBRPC_TEST -LJO "$WEBRPC_TEST_URL" && chmod +x $WEBRPC_TEST

    # # Print test RIDL schema from webrpc-test
    # $WEBRPC_TEST -print-schema > ./test@$VERSION.ridl
    # # Generate server/client code from RIDL schema
    # $WEBRPC_GEN -schema=./test@$VERSION.ridl -target=golang -pkg=tests -client -server -out=./test.gen.go

    echo "###"
    echo "### webrpc@$VERSION server <=> generated client"
    echo "###"
    echo

    # Run webrpc-test@VERSION server and wait for it to be ready
    echo "Running webrpc-test@$VERSION server at 0.0.0.0:$PORT"
    $WEBRPC_TEST -server -port=$PORT -timeout=2s &
    until nc -z localhost $PORT; do sleep 0.1; done

    # Run generated client tests
    go test -v -server=false -client=true -url=http://localhost:$PORT

    wait
    echo

    echo "###"
    echo "### generated server <=> webrpc@$VERSION client"
    echo "###"
    echo

    # Run generated server
    go test -v -server=true -client=false -httptest.serve=0.0.0.0:$PORT -serverTimeout=2s &
    until nc -z localhost $PORT; do sleep 0.1; done

    # Run webrpc-test@VERSION client and wait for it to be ready
    echo "Running webrpc-test@$VERSION client tests against http://localhost:$PORT"
    $WEBRPC_TEST -client -url=http://localhost:$PORT

    wait
    echo
done
