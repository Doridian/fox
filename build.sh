#!/bin/bash
set -e

echo 'Sourcing PKGBUILD'
source ./PKGBUILD
echo 'Sourced.'

set -x

startdir="$(pwd)"
srcdir="${startdir}/src"
mkdir -p "${srcdir}"
prepare
build
