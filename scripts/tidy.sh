#!/bin/bash
# Check go.mod file is up to date

# Check required commands are in place
command -v go >/dev/null 2>&1 || { echo "please install go"; exit 1; }

backup_go_mod()
{
    mod=$(mktemp)
    cp go.mod "$mod"

    sum=$(mktemp)
    cp go.sum "$sum"
}

restore_go_mod()
{
    cp "$mod" go.mod
    rm "$mod"

    cp "$sum" go.sum
    rm "$sum"
}

# Backup actual go.mod and go.sum
backup_go_mod
trap restore_go_mod EXIT

go mod tidy

diff "$mod" go.mod || { echo "go.mod file is not up to date"; exit 42; }
