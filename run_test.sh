#!/bin/bash
set -eo pipefail
cd "$(dirname "$0")"
export PATH=$PATH:/usr/local/go/bin
go test ./... -count=1
