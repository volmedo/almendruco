#!/bin/bash
# Update lambda by uploading a new packaged binary
# Requires an aws cli credentials profile with enough rights to perform the action
set -eufCo pipefail
export SHELLOPTS

# Check required commands are in place
command -v aws >/dev/null 2>&1 || { echo "please install aws"; exit 1; }

usage() {
    echo "usage: $(basename $0) [-r] <function_name> <path_to_package>" >&2
}

if [ $# -lt 2 ] || [ $# -gt 3 ]; then
    echo "[error]: wrong number of parameters provided"
    usage
    exit 1
fi

for_real=false

while getopts 'r' opt; do
    case $opt in
        r) for_real=true ;;
        *) echo '[error]: command line parsing failed' >&2
            exit 1
    esac
done
shift "$(( OPTIND - 1 ))"

function_name="${1}"
path_to_package="${2}"

if [ "$for_real" = true ] ; then
    aws lambda update-function-code \
                --function-name "${function_name}" \
                --zip-file "fileb://${path_to_package}"
else
    aws lambda update-function-code \
                --function-name "${function_name}" \
                --zip-file "fileb://${path_to_package}" \
                --dry-run
fi
