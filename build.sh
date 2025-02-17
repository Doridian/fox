#!/bin/bash
set -e

echo 'Sourcing PKGBUILD'
source ./PKGBUILD
echo 'Sourced.'

set -x

srcdir=~/.local/bin
startdir="$(pwd)"
prepare
build
