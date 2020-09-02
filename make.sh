#!/bin/bash
ARG=$1

COMMAND_NAME=png-to-jpeg
LDFLAGS=-ldflags=-X=main.version=$(git describe)

OS_ARCHS=(
    darwin/386
    darwin/amd64
    linux/386
    linux/amd64
    linux/arm
    linux/arm64
    windows/386
    windows/amd64
)

if [ "$ARG" == "" ]; then
  echo "invalid argument try one the following:"
  echo "  clean       Remove release binaries"
  echo "  release     Build release binaries"
fi

if [ "$ARG" == "release" ]; then
  for OS_ARCH in "${OS_ARCHS[@]}"
  do
    goos=${OS_ARCH%/*}
    goarch=${OS_ARCH#*/}
    GOOS=${goos} GOARCH=${goarch} go build -o "build/$goos/$goarch/$COMMAND_NAME" $LDFLAGS
  done
fi

if [ "$ARG" == "clean" ]; then
  rm -rf build
fi
