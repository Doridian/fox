# Maintainer: Doridian <git at doridian dot net>

# This should ideally be inside a pkgver() subroutine, but that is not possible
# as part of the version comes from the commit count since the latest tag
# so if you commit your current changes the PKGBUILD that would push it one tag further
# than it just calculated, so it would cause a perma-diff in git which is very suboptimal
latest_tag="$(git describe --tags --abbrev=0)"
commits_since_tag="$(git rev-list --count ${latest_tag}..HEAD)"
tag_suffix=""
if ! git status --porcelain; then
  tag_suffix="dev"
fi

pkgname=fox
pkgver="${latest_tag}.${commits_since_tag}${tag_suffix}"
pkgrel="1"
pkgdesc='Fully OwO eXtensions'
arch=('any')
url='https://github.com/Doridian/fox.git'
license=('GPL-3.0-or-later')
makedepends=('git' 'go')
source=()
sha256sums=()

prepare() {
  cd "${startdir}"
  go generate ./...
}

build() {
  cd "${startdir}"
  pkgverfull="${pkgver}-${pkgrel}"
  go build -trimpath -ldflags "-X github.com/Doridian/fox/modules/info.version=${pkgverfull} -X github.com/Doridian/fox/modules/info.gitrev=$(git rev-parse HEAD)" -o "${srcdir}/fox" ./cmd
}

package() {
  cd "${srcdir}"
  install -Dm755 ./fox "${pkgdir}/bin/fox"
}
