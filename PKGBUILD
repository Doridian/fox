# Maintainer: Doridian <git at doridian dot net>

latest_tag="$(git describe --tags --abbrev=0)"
commits_since_tag="$(git rev-list --count ${latest_tag}..HEAD)"
tag_suffix=""
if ! git diff-index --quiet HEAD --; then
  tag_suffix="dev"
fi

pkgname=fox
pkgver="${latest_tag}${tag_suffix}"
pkgrel="$((commits_since_tag + 1))"
pkgdesc='Fully OwO eXtensions'
arch=('any')
url='https://github.com/Doridian/fox.git'
license=('GPL-3.0')
makedepends=('git' 'go')
source=()
sha256sums=()

prepare() {
  cd "${srcdir}"
  go generate ../...
}

build() {
  cd "${srcdir}"
  pkgverfull="${pkgver}-${pkgrel}"
  go build -trimpath -ldflags "-X github.com/Doridian/fox/modules/info.version=${pkgverfull} -X github.com/Doridian/fox/modules/info.gitrev=$(git rev-parse HEAD)" -o ./fox ../cmd
}

package() {
  cd "${srcdir}"
  install -Dm755 ./fox "${pkgdir}/bin/fox"
}

# vim:set ts=2 sw=2 et:

