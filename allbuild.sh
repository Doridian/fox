#!/bin/bash
set -e

echo 'source ./PKGBUILD'
source ./PKGBUILD
echo './PKGBUILD OK!'
set -x

startdir="$(pwd)"
mkdir -p pkg src

pkgdir="${startdir}/pkg"
srcdir="${startdir}/src"

prepare

dobuild() {
    build

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
    mv "${srcdir}/fox" "${pkgdir}/fox${OSSUFFIX}${ARCHSUFFIX}${XSUFFIX}"
}

GOOS=linux GOARCH=amd64 dobuild
GOOS=linux GOARCH=arm64 dobuild
GOOS=windows GOARCH=amd64 dobuild
GOOS=windows GOARCH=arm64 dobuild
GOOS=darwin GOARCH=arm64 dobuild
