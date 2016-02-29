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
REGISTRY ?= johnt337

# Get the git commit
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)

# ECR tagging
ECR_TAG=latest
ECR_VERSION=$(shell grep -E 'Version =' ./cmd/ecr_login/version.go | awk '{print$$NF}' | sed 's@"@@g')

build: godeps build-all

build-all:
	@make build-awscli
	@make build-ecr_login

build-awscli:
	@echo "running make build-awscli"
	docker build -t awscli-build -f Dockerfile .
	GO15VENDOREXPERIMENT=$(GO15VENDOREXPERIMENT) docker run --rm -v /var/run:/var/run -v $(MOUNT):/go/src/$(PROJ) --entrypoint=/bin/sh -i awscli-build -c "godep restore && make lint && make lint-check && make test/units && make docker/awscli "

build-ecr_login:
	@echo "running make build-ecr_login"
	docker build -t awscli-build -f Dockerfile .
	GO15VENDOREXPERIMENT=$(GO15VENDOREXPERIMENT) docker run --rm -v /var/run:/var/run -v $(MOUNT):/go/src/$(PROJ) --entrypoint=/bin/sh -i awscli-build -c "godep restore && make lint && make lint-check && make test/units && make docker/ecr_login "


docker/awscli: $(SRC) config Dockerfile.awscli
	@echo "running make docker/awscli"
	make bin/awscli
	[ -d ./tmp ] || mkdir ./tmp && chmod 4777 ./tmp
	docker build -t $(REGISTRY)/awscli -f Dockerfile.awscli .

docker/ecr_login: $(SRC) config Dockerfile.ecr_login
	@echo "running make docker/ecr_login"
	make bin/ecr_login
	[ -d ./tmp ] || mkdir ./tmp && chmod 4777 ./tmp
	[ -d ./certs ] || cp -a /etc/ssl/certs .
	docker build -t $(REGISTRY)/ecr_login:$(ECR_VERSION) -f Dockerfile.ecr_login .
	docker build -t $(REGISTRY)/ecr_login:$(ECR_VERSION)-docker -f ecr_login_plus_docker.Dockerfile .
	docker tag -f $(REGISTRY)/ecr_login:$(ECR_VERSION) $(REGISTRY)/ecr_login:$(ECR_TAG)


bin: $(SRC)
	@make bin/awscli
	@make bin/ecr_login

bin/awscli: $(SRC)
	@echo "statically linking awscli"
	CGO_ENABLED=0 GOOS=linux godep go build -a -installsuffix cgo -ldflags '-w -X main.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)' -o bin/awscli cmd/awscli/*.go

bin/ecr_login: $(SRC)
	@echo "statically linking ecr_login"
	CGO_ENABLED=0 GOOS=linux godep go build -a -installsuffix cgo -ldflags '-w -X main.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)' -o bin/ecr_login cmd/ecr_login/*.go

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
	rm -f bin/awscli bin/ecr_login
	docker images | grep -E '<none>' | awk '{print$$3}' | xargs docker rmi

distclean:
	@make clean
	@echo "running make distclean"
	rm -rf ./tmp ./certs
	docker rm awscli-build run-awscli
	docker rmi bin/awscli bin/ecr_login awscli-go $(REGISTRY)/awscli $(REGISTRY)/ecr_login

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

run-ecr_login: config
	@echo "running bin/ecr_login"
	docker run -it --rm -i bin/ecr_login

test:
	@echo "running test"
	docker run -it --rm -v /var/run:/var/run -v $(MOUNT):/go/src/$(PROJ) --entrypoint=/bin/sh -i awscli-build -c "godep restore && make test/units"

test/units: $(SRC)
	@echo "running test/units"
	godep go test $(TESTS) $(ARGS)

test-cover: $(SRC)
	@godep go test $(COVER) -coverprofile=coverage.ou
	@godep go tool cover -html=coverage.out

.PHONY: clean test test/units run-bin/awscli run-bin/ecr_login interactive bootstrap bootstrap-test lint lint-check test-cover godeps
