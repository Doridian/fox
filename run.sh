#!/bin/bash
set -e

echo 'source ./PKGBUILD'
source ./PKGBUILD
echo './PKGBUILD OK!'
set -x

goldflags="${GOLDFLAGS-}"

startdir="$(pwd)"
mkdir -p src
srcdir="${startdir}/src"

prepare
build

if [-z "${BINPATH-}" ]; then
    exec ./src/fox
else
    mv ./src/fox "${BINPATH}"
    exec "${BINPATH}"
fi
