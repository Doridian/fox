# Maintainer: Doridian <git at doridian dot net>

pkgname=fox
pkgver=0.1.0
pkgrel=1
pkgdesc='Fully OwO eXtensions'
arch=('any')
url='https://github.com/Doridian/fox.git'
license=('GPL-3.0')
makedepends=('go')
source=()
sha256sums=()

build() {
  cd "${srcdir}"
  go build -trimpath -o ./fox ../cmd
}

package() {
  cd "${srcdir}"
  install -Dm755 ./fox "${pkgdir}/bin/fox"
}

# vim:set ts=2 sw=2 et:

