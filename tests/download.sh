#!/bin/bash
set -e

VERSION="${1}"
DIR="${2}"
[[ -z "$VERSION" || -z "$DIR" ]] && { echo "Usage: $0 <webrpc-version> <dir>"; exit 1; }

mkdir -p "$DIR"

# Download webrpc binaries if not available locally
OS="$(basename $(uname -o | tr A-Z a-z))"
ARCH="$(uname -m | sed 's/x86_64/amd64/')"
if [[ ! -f "$DIR/webrpc-gen" ]]; then
    curl -o "$DIR/webrpc-gen" -fLJO "https://github.com/webrpc/webrpc/releases/download/$VERSION/webrpc-gen.$OS-$ARCH"
    chmod +x "$DIR/webrpc-gen"
fi
if [[ ! -f "$DIR/webrpc-test" ]]; then
    curl -o "$DIR/webrpc-test" -fLJO "https://github.com/webrpc/webrpc/releases/download/$VERSION/webrpc-test.$OS-$ARCH"
    chmod +x "$DIR/webrpc-test"
fi

