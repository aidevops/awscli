#!/bin/sh

. $PWD/go_path.env
echo "GOPATH=$GOPATH"
echo "PATH=$PATH"
echo "GO15VENDOREXPERIMENT=$GO15VENDOREXPERIMENT"

echo "installing go if not already installed..."
sudo apk add --update go && \
. $PWD/go_get.sh
