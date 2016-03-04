#!/bin/sh

. $PWD/go_path.env
echo "GOPATH=$GOPATH"
echo "PATH=$PATH"
echo "GO15VENDOREXPERIMENT=$GO15VENDOREXPERIMENT"

echo "installing go dev deps"
go get github.com/tools/godep && \
go get github.com/onsi/ginkgo/ginkgo && \
go get github.com/onsi/gomega && \
go get -u github.com/golang/lint/golint && \
go get golang.org/x/tools/cmd/cover && \
echo "installing build deps..." && \
go get -u github.com/aws/aws-sdk-go && \
go get github.com/mitchellh/cli && \
go get github.com/fatih/color && \
go get github.com/go-ini/ini && \
go get github.com/jmespath/go-jmespath && \
go get github.com/mattn/go-colorable && \
go get golang.org/x/sys/unix && \
echo "updating godep vendoring" && \
godep save ./...
