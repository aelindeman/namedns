#!/usr/bin/env bash

set -e

mkdir -p dist

if [ "$CI" == "true" ]; then
  VERSION="${TRAVIS_TAG}"
  COMMIT_HASH="$TRAVIS_COMMIT; $TRAVIS_BUILD_ID"
else
  VERSION=$(git tag -l --points-at HEAD)
  COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null)
fi

[ -z "$VERSION" ] && VERSION="$(git rev-parse --abbrev-ref HEAD 2>/dev/null)"

GO_BUILD_CMD="go build"
GO_BUILD_LDFLAGS="-X github.com/aelindeman/namedns/cmd.Version=${VERSION}+${COMMIT_HASH}"

[ -z "$GO_BUILD_PLATFORMS" ] && GO_BUILD_PLATFORMS="linux windows darwin"
[ -z "$GO_BUILD_ARCHS" ] && GO_BUILD_ARCHS="amd64"

for OS in ${GO_BUILD_PLATFORMS[@]}; do
  for ARCH in ${GO_BUILD_ARCHS[@]}; do
    NAME="namedns-$OS-$ARCH"
    [ "$OS" == "windows" ] && NAME="$NAME.exe"
    echo "building $NAME"
    GOARCH=$ARCH GOOS=$OS $GO_BUILD_CMD -ldflags "$GO_BUILD_LDFLAGS" -o "dist/$NAME"
    echo "checksumming $NAME"
    shasum -a 256 "dist/$NAME" > "dist/$NAME".sha256sum
  done
done
