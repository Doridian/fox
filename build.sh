#!/bin/sh
set -ex
go generate ./...
go build -trimpath -o ~/.local/bin/fox ./cmd/fox.go
exec fox
