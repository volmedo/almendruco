#!/bin/bash
# Lint code
set -eufCo pipefail
export SHELLOPTS

# Check required commands are in place
command -v golangci-lint >/dev/null 2>&1 || { echo "please install golangci-lint"; exit 1; }

golangci-lint run -E goimports
