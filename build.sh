#!/bin/sh
set -ex
go generate ./...
go build -trimpath -o ~/.local/bin/fox ./cmd
exec fox
