#!/bin/bash
# Package a binary into file suitable to be deployed
set -eufCo pipefail
export SHELLOPTS

usage() {
    echo "usage: $(basename $0) <bin_name> <package_name>" >&2
}

if [ $# -ne 2 ] ; then
    echo "[error]: both a binary name and a package name must be provided"
    usage
    exit 1
fi

bin_name="${1}"
package_name="${2}"

zip ${package_name} ${bin_name}
