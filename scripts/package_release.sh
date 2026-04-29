#!/usr/bin/env bash
set -euo pipefail

APP_NAME="pokedex-fetch-go"
CMD_PATH="./cmd/pokedex-fetch-go"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

usage() {
  cat <<EOF
Usage: ./scripts/package_release.sh vX.Y.Z

Cross-compiles release binaries, packages them with the executable named
'pokedex-fetch-go' inside each archive, and writes releasable assets to:

  dist/vX.Y.Z/

Example:

  ./scripts/package_release.sh v1.0.0
EOF
}

if [[ $# -ne 1 ]]; then
  usage >&2
  exit 1
fi

VERSION="$1"
if [[ "$VERSION" != v* ]]; then
  VERSION="v$VERSION"
fi

if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+([-+][0-9A-Za-z.-]+)?$ ]]; then
  echo "invalid version: $1" >&2
  echo "expected format like v1.0.0" >&2
  exit 1
fi

for tool in go tar sha256sum zip; do
  if ! command -v "$tool" >/dev/null 2>&1; then
    echo "required tool not found: $tool" >&2
    exit 1
  fi
done

cd "$ROOT"

OUT_DIR="$ROOT/dist/$VERSION"
WORK_DIR="$ROOT/dist/.work-$VERSION"
PKG_DIR="$WORK_DIR/pkg"

rm -rf "$OUT_DIR" "$WORK_DIR"
mkdir -p "$OUT_DIR" "$PKG_DIR"

build_and_package() {
  local goos="$1"
  local goarch="$2"
  local ext="${3:-}"
  local archive_ext="$4"
  local asset_base="${APP_NAME}_${VERSION}_${goos}_${goarch}"
  local pkg_path="$PKG_DIR/$asset_base"
  local bin_name="$APP_NAME$ext"

  echo "==> Building $goos/$goarch"
  mkdir -p "$pkg_path"
  GOOS="$goos" GOARCH="$goarch" go build -trimpath -ldflags "-s -w" -o "$pkg_path/$bin_name" "$CMD_PATH"

  cp LICENSE "$pkg_path/"
  cp README.md "$pkg_path/"
  cp config.example.toml "$pkg_path/"

  if [[ "$archive_ext" == "zip" ]]; then
    echo "==> Packaging $asset_base.zip"
    (cd "$PKG_DIR" && zip -qr "$OUT_DIR/$asset_base.zip" "$asset_base")
  else
    echo "==> Packaging $asset_base.tar.gz"
    tar -czf "$OUT_DIR/$asset_base.tar.gz" -C "$PKG_DIR" "$asset_base"
  fi
}

build_and_package linux amd64 "" tar.gz
build_and_package linux arm64 "" tar.gz
build_and_package darwin amd64 "" tar.gz
build_and_package darwin arm64 "" tar.gz
build_and_package windows amd64 ".exe" zip

(
  cd "$OUT_DIR"
  sha256sum * > checksums.txt
)

rm -rf "$WORK_DIR"

echo
echo "Release assets ready: $OUT_DIR"
ls -lh "$OUT_DIR"
