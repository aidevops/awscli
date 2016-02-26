#!/bin/sh

export GOPATH=$PWD/../../../../
export PATH=$GOPATH/bin:$PATH

sudo apk add --update go && \
  go get github.com/tools/godep && \
  go get -u github.com/aws/aws-sdk-go && \
  godep save -r
