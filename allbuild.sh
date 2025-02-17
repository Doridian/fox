#!/bin/bash
set -e

source ./PKGBUILD

startdir="$(pwd)"
mkdir -p pkg src

srcdir="${startdir}/src"

prepare
cd "${startdir}"

dobuild() {
    build
    cd "${startdir}"
    XSUFFIX=""
    OSSUFFIX="-${GOOS}"
    ARCHSUFFIX="-${GOARCH}"
    if [ "${GOOS}" == "windows" ]; then
        XSUFFIX=".exe"
        OSSUFFIX=""
        if [ "${GOARCH}" == "amd64" ]; then
            ARCHSUFFIX=""
        fi
    elif [ "${GOOS}" == "darwin" ]; then
        OSSUFFIX="-macos"
        if [ "${GOARCH}" == "arm64" ]; then
            ARCHSUFFIX=""
        fi
    fi
    mv src/fox "pkg/fox${OSSUFFIX}${ARCHSUFFIX}${XSUFFIX}"
}

GOOS=linux GOARCH=amd64 dobuild
GOOS=linux GOARCH=arm64 dobuild
GOOS=windows GOARCH=amd64 dobuild
GOOS=windows GOARCH=arm64 dobuild
GOOS=darwin GOARCH=arm64 dobuild
