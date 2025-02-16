#!/bin/sh
set -e

echo 'Sourcing PKGBUILD'
source ./PKGBUILD
srcdir=~/.local/bin
startdir="$(pwd)"
echo 'prepare()'
prepare
echo 'build()'
build
echo 'Done, exec...'

exec fox
