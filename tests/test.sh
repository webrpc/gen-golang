#!/bin/bash
set -e

### Run interoperability tests against specific version of webrpc-test binary.

PORT=9889
VERSION="$(webrpc-test -version)"

echo "###"
echo "### reference server ($VERSION) <=> generated client"
echo "###"
echo

# Run reference server and wait for it to be ready
echo "Running reference server at 0.0.0.0:$PORT"
webrpc-test -server -port=$PORT -timeout=2s &

# Wait until http://localhost:$PORT is available, up to 10s.
for (( i=0; i<100; i++ )); do nc -z localhost $PORT && break || sleep 0.1; done

# Run generated client tests
go test -v -server=false -client=true -url=http://localhost:$PORT

wait
echo

echo "###"
echo "### generated server <=> reference client ($VERSION)"
echo "###"
echo

# Run generated server
go test -v -server=true -client=false -httptest.serve=0.0.0.0:$PORT -serverTimeout=2s &

# Wait until http://localhost:$PORT is available, up to 10s.
for (( i=0; i<100; i++ )); do nc -z localhost $PORT && break || sleep 0.1; done

# Run reference client and wait for it to be ready
echo "Running reference client tests against http://localhost:$PORT"
webrpc-test -client -url=http://localhost:$PORT

wait
echo
