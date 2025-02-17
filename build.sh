#!/bin/bash
set -e

source ./PKGBUILD

startdir="$(pwd)"
srcdir="${startdir}/src"
mkdir -p "${srcdir}"

prepare
build
