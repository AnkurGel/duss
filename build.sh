#!/bin/bash
export PATH="$PATH:$GOBIN"
cd $GOPATH/src/github.com/ankurgel/duss
echo "Installing build dependencies..."
go get -v github.com/mitchellh/gox
go get -v github.com/gobuffalo/packr/v2/packr2
echo "Installing dependencies..."
dep ensure
echo "Packing build..."
cd ./cmd/duss && packr2 clean && packr2 && cd ../../
rm -rf releases
mkdir releases
echo "Preparing binaries..."
gox -output="releases/{{.Dir}}_{{.OS}}_{{.Arch}}" -osarch="linux/amd64 darwin/amd64" ./cmd/duss
echo "Cleaning up..."
cd ./cmd/duss && packr2 clean && cd ../../
echo "Binaries ready in ./releases/"
