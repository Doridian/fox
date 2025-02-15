# Maintainer: Doridian <git at doridian dot net>

pkgname=fox-git
pkgver=r194.ed4bb93
pkgrel=1
pkgdesc='Fully OwO eXtensions'
arch=('any')
url='https://github.com/Doridian/fox.git'
license=('GPL-3.0')
makedepends=('git' 'go')
source=()
sha256sums=()

pkgver() {
  printf "r%s.%s" "$(git rev-list --count HEAD)" "$(git rev-parse --short HEAD)"
}

build() {
  cd "${srcdir}"
  go build -trimpath -o ./fox ../cmd
}

package() {
  cd "${srcdir}"
  install -Dm755 ./fox "${pkgdir}/bin/fox"
}

# vim:set ts=2 sw=2 et:

