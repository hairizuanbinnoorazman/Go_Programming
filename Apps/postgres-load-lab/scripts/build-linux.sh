#!/usr/bin/env bash
set -euo pipefail

project_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
mkdir -p "$project_dir/bin"
cd "$project_dir"
GOWORK=off CGO_ENABLED=0 GOOS=linux GOARCH="${GOARCH:-amd64}" go build -trimpath -o bin/loadlab-server ./cmd/server
GOWORK=off CGO_ENABLED=0 GOOS=linux GOARCH="${GOARCH:-amd64}" go build -trimpath -o bin/loadlab-loadgen ./cmd/loadgen
