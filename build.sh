#!/bin/sh
set -ex

source ./PKGBUILD
srcdir=~/.local/bin
prepare
build

exec fox
