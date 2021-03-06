#!/bin/bash
export PATH="$PATH:$GOBIN"

export BUILD_TYPE=${1:-"debug"}
export LDFLAGS=''
if [[ $BUILD_TYPE == 'release' ]]; then
  export LDFLAGS='-s -w'
fi
echo "Building for $BUILD_TYPE"

cd $GOPATH/src/github.com/ankurgel/duss
echo "Installing build dependencies..."
go get -v github.com/mitchellh/gox
go get -v github.com/gobuffalo/packr/v2/packr2
echo "Installing dependencies..."
dep ensure
echo "Packing build..."
cd ./cmd/duss && packr2 clean && packr2
rm -rf ../../releases
mkdir ../../releases
echo "Preparing binaries..."
gox -output="../../releases/{{.Dir}}_{{.OS}}_{{.Arch}}" -osarch="linux/amd64 darwin/amd64" -ldflags="$LDFLAGS"
echo "Cleaning up..."
packr2 clean
cd ../../
echo "Binaries ready in ./releases/"
