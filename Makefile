MAKEFLAGS  += --no-builtin-rules
.SUFFIXES:
.SECONDARY:
.DELETE_ON_ERROR:

ARGS  ?= -v -race
PROJ  ?= github.com/johnt337/awscli
MAIN  ?= $(go list ./... | grep -v /vendor/)
TESTS ?= $(MAIN) -cover
LINTS ?= $(MAIN)
COVER ?=
SRC   := $(shell find . -name '*.go')
MOUNT ?= $(shell pwd)
GO15VENDOREXPERIMENT ?= 1

# Get the git commit
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)

build: godeps build-awscli

build-awscli:
	@echo "running make build-awscli"
	docker build -t awscli-build -f Dockerfile .
	GO15VENDOREXPERIMENT=$(GO15VENDOREXPERIMENT) docker run --rm -v /var/run:/var/run -v $(MOUNT):/go/src/$(PROJ) --entrypoint=/bin/sh -i awscli-build -c "godep restore && make lint && make lint-check && make test/awscli && make docker/awscli "

docker/awscli: $(SRC) config Dockerfile.awscli
	@echo "running make docker/awscli"
	make bin/awscli
	[ -d ./tmp ] || mkdir ./tmp && chmod 4777 ./tmp
	docker build -t bin/awscli -f Dockerfile.awscli .

bin/awscli: $(SRC)
	@echo "statically linking awscli"
	CGO_ENABLED=0 GOOS=linux godep go build -a -installsuffix cgo -ldflags '-w -X main.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)' -o bin/awscli cmd/awscli/*.go

bootstrap:
	ginkgo bootstrap

bootstrap-test:
	@cd ${DIR} && ginkgo bootstrap && cd ~-

godeps:
	@echo "running godep"
	godep save ./...
	#godep go build ./...
	#godep go test $(MAIN)

clean:
	@echo "running make clean"
	rm -f bin/awscli
	docker images | grep -E '<none>' | awk '{print$$3}' | xargs docker rmi

distclean:
	@make clean
	@echo "running make distclean"
	rm -rf ./tmp ./certs
	docker rm awscli-build run-awscli
	docker rmi bin/awscli awscli-go

interactive:
	@echo "running make build-awscli"
	docker build -t awscli-go -f Dockerfile .
	docker run -it --rm --name awscli-build -v /var/run:/var/run -v $(MOUNT):/go/src/$(PROJ) --entrypoint=/bin/bash -i awscli-build

lint: $(SRC)
	@golint $(LINTS)

lint-check: $(SRC)
	@echo "`golint $(LINTS) | wc -l | awk '{print$$NF}'` error(s)"
	@[ `golint $(LINTS) | wc -l | awk '{print$$NF}'` -le 0 ] && true || false

run-awscli: config
	@echo "running bin/awscli"
	docker run -it --rm -i bin/awscli

test:
	@echo "running test"
	docker run -it --rm -v /var/run:/var/run -v $(MOUNT):/go/src/$(PROJ) --entrypoint=/bin/sh -i awscli-build -c "godep restore && make test/awscli"

test/awscli: $(SRC)
	@echo "running test/awscli"
	godep go test $(TESTS) $(ARGS)

test-cover: $(SRC)
	@godep go test $(COVER) -coverprofile=coverage.ou
	@godep go tool cover -html=coverage.out

.PHONY: clean test test/awscli run-bin/awscli interactive bootstrap bootstrap-test lint lint-check test-cover godeps
