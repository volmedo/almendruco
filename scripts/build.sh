#!/bin/bash
# Build almendruco binary from sources
set -eufCo pipefail
export SHELLOPTS

# Check required commands are in place
command -v go >/dev/null 2>&1 || { echo "please install go"; exit 1; }

usage() {
    echo "usage: $(basename $0) <bin_name>" >&2
}

if [ $# -ne 1 ] ; then
    echo "[error]: please provide a name for the resulting binary"
    usage
    exit 1
fi

bin_name="${1}"

go build -ldflags "-s -w" -o "${bin_name}" "${PWD}/cmd/almendruco"
