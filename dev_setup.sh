#!/bin/sh

export GOPATH=$PWD/../../../
export PATH=$GOPATH/bin:$PATH

sudo apk add --update go && \
  go get github.com/tools/godep && \
  godep save -r
