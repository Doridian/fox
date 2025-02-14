#!/bin/sh
set -ex
go generate ./...
go build -o ~/.local/bin/fox ./cmd
exec fox
