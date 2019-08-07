#!/usr/bin/env bash

set -e

if [[ $(uname) = "Linux" ]] && which apt &> /dev/null; then
  apt update
  apt install libpcap-dev -y
fi

if [[ $(uname) = "Linux" ]] && which apk &> /dev/null; then
  apk add libpcap-dev git gcc libc-dev
fi

PROJECT=$(basename $(git rev-parse --show-toplevel))
VERSION=$(cat version.txt)
# source .envrc

echo "Running: go test ./..."
go test ./...

GOARCH=amd64

if [[ $(uname) = "Darwin" ]]; then
  GOOS=darwin

  echo "go build -o bin/${PROJECT}-${VERSION}-${GOOS}-${GOARCH}"
  go build -o "bin/${PROJECT}-${VERSION}-${GOOS}-${GOARCH}"
  shasum -a 256 < "bin/${PROJECT}-${VERSION}-${GOOS}-${GOARCH}" | cut -d ' ' -f1

  # docker run -it --rm -v $(pwd):/$(pwd) -w $(pwd) golang:alpine sh ./build.sh # Alpine
  docker run -it --rm -v $(pwd):/$(pwd) -w $(pwd) golang ./build.sh # Ubuntu
fi

if [[ $(uname) = "Linux" ]]; then
  GOOS=linux

  echo "go build -o bin/${PROJECT}-${VERSION}-${GOOS}-${GOARCH}"
  go build -o "bin/${PROJECT}-${VERSION}-${GOOS}-${GOARCH}"
  shasum -a 256 < "bin/${PROJECT}-${VERSION}-${GOOS}-${GOARCH}" | cut -d ' ' -f1
fi



