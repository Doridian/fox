#!/bin/sh
set -e
go build -o ~/.local/bin/fox .
exec fox
