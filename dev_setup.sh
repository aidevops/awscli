#!/bin/sh

export GOPATH=$PWD/../../../../
export PATH=$GOPATH/bin:$PATH
export GO15VENDOREXPERIMENT=1

sudo apk add --update go && \
  go get github.com/tools/godep && \
  go get github.com/onsi/ginkgo/ginkgo && \
  go get github.com/onsi/gomega && \
  go get -u github.com/golang/lint/golint && \
  go get -u github.com/aws/aws-sdk-go && \
  go get golang.org/x/tools/cmd/cover && \
  go get github.com/mitchellh/cli && \
  go get golang.org/x/sys/unix && \
  godep save ./...
