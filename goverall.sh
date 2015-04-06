#!/bin/bash
set -e

mkdir -p .cover
go list ./... | xargs -I% bash -c 'name="%"; go test % --coverprofile=.cover/${name//\//_}'
echo "mode: set" > cover.out
cat .cover/* | grep -v mode >> cover.out
rm -r .cover

if [[ "$1" != "" ]]; then
	go tool cover -$1=cover.out
fi
