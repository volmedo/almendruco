#!/bin/bash
# Run unit tests and get coverage
set -eufCo pipefail
export SHELLOPTS

# Check required commands are in place
command -v go >/dev/null 2>&1 || { echo "please install go"; exit 1; }

go test -race -coverprofile=.test_coverage.txt ./...
go tool cover -func=.test_coverage.txt | tail -n1 | awk '{print "Total test coverage: " $3}'
