MAKEFLAGS  += --no-builtin-rules
.SUFFIXES:
.SECONDARY:
.DELETE_ON_ERROR:

ARGS  ?= -v -race
PROJ  ?= github.com/aidevops/awscli
MAIN  ?= $(shell go list ./... | grep -v /vendor/)
TESTS ?= $(MAIN) -cover
LINTS ?= $(MAIN)
COVER ?=
SRC   := $(shell find . -name '*.go')
MOUNT ?= $(shell pwd)
REGISTRY ?= aidevops

# Get the git commit
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)

# AWSCLI tagging
AWSCLI_TAG=latest
AWSCLI_VERSION=$(shell grep -E 'Version =' ./cmd/awscli/version.go | awk '{print$$NF}' | sed 's@"@@g')

# ECR tagging
ECR_TAG=latest
ECR_VERSION=$(shell grep -E 'Version =' ./cmd/ecr_login/version.go | awk '{print$$NF}' | sed 's@"@@g')

# EC2 tagging
EC2_TAG=latest
EC2_TAG_VERSION=$(shell grep -E 'Version =' ./cmd/ec2_tag/version.go | awk '{print$$NF}' | sed 's@"@@g')

# SQS messaging
SQS_TAG=latest
SQS_VERSION=$(shell grep -E 'Version =' ./cmd/sqs_util/version.go | awk '{print$$NF}' | sed 's@"@@g')

# S3 storage
S3_TAG=latest
S3_VERSION=$(shell grep -E 'Version =' ./cmd/s3_util/version.go | awk '{print$$NF}' | sed 's@"@@g')

# SG registration
SG_REGISTER_TAG=latest
SG_REGISTER_VERSION=$(shell grep -E 'Version =' ./cmd/sg_register/version.go | awk '{print$$NF}' | sed 's@"@@g')

build: build-all

build-all:
	@make build-awscli
	@make build-ecr_login
	@make build-ec2_tag
	@make build-sqs_util
	@make build-s3_util
	@make build-sg_register

build-awscli:
	@echo "running make build-awscli"
	docker build -t awscli-build -f Dockerfile .
	docker run --rm -v /var/run:/var/run -v $(MOUNT):/go/src/$(PROJ) --entrypoint=/bin/sh -i awscli-build -c "godep restore && make lint && make lint-check && make test/units && make docker/awscli "

build-ecr_login:
	@echo "running make build-ecr_login"
	docker build -t awscli-build -f Dockerfile .
	docker run --rm -v /var/run:/var/run -v $(MOUNT):/go/src/$(PROJ) --entrypoint=/bin/sh -i awscli-build -c "godep restore && make lint && make lint-check && make test/units && make docker/ecr_login "

build-ec2_tag:
	@echo "running make build-ec2_tag"
	docker build -t awscli-build -f Dockerfile .
	docker run --rm -v /var/run:/var/run -v $(MOUNT):/go/src/$(PROJ) --entrypoint=/bin/sh -i awscli-build -c "godep restore && make lint && make lint-check && make test/units && make docker/ec2_tag "

build-sqs_util:
	@echo "running make build-sqs_util"
	docker build -t awscli-build -f Dockerfile .
	docker run --rm -v /var/run:/var/run -v $(MOUNT):/go/src/$(PROJ) --entrypoint=/bin/sh -i awscli-build -c "godep restore && make lint && make lint-check && make test/units && make docker/sqs_util "

build-s3_util:
	@echo "running make build-s3_util"
	docker build -t awscli-build -f Dockerfile .
	docker run --rm -v /var/run:/var/run -v $(MOUNT):/go/src/$(PROJ) --entrypoint=/bin/sh -i awscli-build -c "godep restore && make lint && make lint-check && make test/units && make docker/s3_util "

build-sg_register:
	@echo "running make build-sg_register"
	docker build -t awscli-build -f Dockerfile .
	docker run --rm -v /var/run:/var/run -v $(MOUNT):/go/src/$(PROJ) --entrypoint=/bin/sh -i awscli-build -c "godep restore && make lint && make lint-check && make test/units && make docker/sg_register "

docker/awscli: $(SRC) awscli.Dockerfile
	@echo "running make docker/awscli"
	make bin/awscli
	[ -d ./tmp ] || mkdir ./tmp && chmod 4777 ./tmp
	docker build -t $(REGISTRY)/awscli:$(AWSCLI_VERSION) -f awscli.Dockerfile .
	docker tag $(REGISTRY)/awscli:$(AWSCLI_VERSION) $(REGISTRY)/awscli:$(AWSCLI_TAG)

docker/ecr_login: $(SRC) ecr_login.Dockerfile
	@echo "running make docker/ecr_login"
	make bin/ecr_login
	[ -d ./tmp ] || mkdir ./tmp && chmod 4777 ./tmp
	[ -d ./certs ] || cp -a /etc/ssl/certs .
	docker build -t $(REGISTRY)/ecr_login:$(ECR_VERSION) -f ecr_login.Dockerfile .
	docker build -t $(REGISTRY)/ecr_login:$(ECR_VERSION)-docker -f ecr_login_plus_docker.Dockerfile .
	docker tag $(REGISTRY)/ecr_login:$(ECR_VERSION) $(REGISTRY)/ecr_login:$(ECR_TAG)

docker/ec2_tag: $(SRC) ec2_tag.Dockerfile
	@echo "running make docker/ec2_tag"
	make bin/ec2_tag
	[ -d ./tmp ] || mkdir ./tmp && chmod 4777 ./tmp
	[ -d ./certs ] || cp -a /etc/ssl/certs .
	docker build -t $(REGISTRY)/ec2_tag:$(EC2_TAG_VERSION) -f ec2_tag.Dockerfile .
	docker tag $(REGISTRY)/ec2_tag:$(EC2_TAG_VERSION) $(REGISTRY)/ec2_tag:$(EC2_TAG)

docker/sqs_util: $(SRC) sqs_util.Dockerfile
	@echo "running make docker/sqs_util"
	make bin/sqs_util
	[ -d ./tmp ] || mkdir ./tmp && chmod 4777 ./tmp
	[ -d ./certs ] || cp -a /etc/ssl/certs .
	docker build -t $(REGISTRY)/sqs_util:$(SQS_VERSION) -f sqs_util.Dockerfile .
	docker tag $(REGISTRY)/sqs_util:$(SQS_VERSION) $(REGISTRY)/sqs_util:$(SQS_TAG)

docker/s3_util: $(SRC) s3_util.Dockerfile
	@echo "running make docker/s3_util"
	make bin/s3_util
	[ -d ./tmp ] || mkdir ./tmp && chmod 4777 ./tmp
	[ -d ./certs ] || cp -a /etc/ssl/certs .
	docker build -t $(REGISTRY)/s3_util:$(S3_VERSION) -f s3_util.Dockerfile .
	docker tag $(REGISTRY)/s3_util:$(S3_VERSION) $(REGISTRY)/s3_util:$(S3_TAG)

docker/sg_register: $(SRC) sg_register.Dockerfile
	@echo "running make docker/sg_register"
	make bin/sg_register
	[ -d ./tmp ] || mkdir ./tmp && chmod 4777 ./tmp
	[ -d ./certs ] || cp -a /etc/ssl/certs .
	docker build -t $(REGISTRY)/sg_register:$(SG_REGISTER_VERSION) -f sg_register.Dockerfile .
	docker tag $(REGISTRY)/sg_register:$(SG_REGISTER_VERSION) $(REGISTRY)/sg_register:$(SG_REGISTER_TAG)

docker/push:
	docker push $(REGISTRY)/awscli:$(AWSCLI_VERSION)
	docker push $(REGISTRY)/awscli:$(AWSCLI_TAG)
	docker push $(REGISTRY)/ecr_login:$(ECR_VERSION)
	docker push $(REGISTRY)/ecr_login:$(ECR_VERSION)-docker
	docker push $(REGISTRY)/ecr_login:$(ECR_TAG)
	docker push $(REGISTRY)/ec2_tag:$(EC2_TAG_VERSION)
	docker push $(REGISTRY)/ec2_tag:$(EC2_TAG)
	docker push $(REGISTRY)/sqs_util:$(SQS_VERSION)
	docker push $(REGISTRY)/sqs_util:$(SQS_TAG)
	docker push $(REGISTRY)/s3_util:$(S3_VERSION)
	docker push $(REGISTRY)/s3_util:$(S3_TAG)
	docker push $(REGISTRY)/sg_register:$(SG_REGISTER_VERSION)
	docker push $(REGISTRY)/sg_register:$(SG_REGISTER_TAG)

bin: $(SRC)
	@make bin/awscli
	@make bin/ecr_login
	@make bin/ec2_tag
	@make bin/sqs_util
	@make bin/s3_util
	@make bin/sg_register

bin/awscli: $(SRC)
	@echo "statically linking awscli"
	CGO_ENABLED=0 GOOS=linux godep go build -a -installsuffix cgo -ldflags '-w -X main.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)' -o bin/awscli cmd/awscli/*.go

bin/ecr_login: $(SRC)
	@echo "statically linking ecr_login"
	CGO_ENABLED=0 GOOS=linux godep go build -a -installsuffix cgo -ldflags '-w -X main.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)' -o bin/ecr_login cmd/ecr_login/*.go

bin/ec2_tag: $(SRC)
	@echo "statically linking ec2_tag"
	CGO_ENABLED=0 GOOS=linux godep go build -a -installsuffix cgo -ldflags '-w -X main.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)' -o bin/ec2_tag cmd/ec2_tag/*.go

bin/sqs_util: $(SRC)
	@echo "statically linking sqs_util"
	CGO_ENABLED=0 GOOS=linux godep go build -a -installsuffix cgo -ldflags '-w -X main.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)' -o bin/sqs_util cmd/sqs_util/*.go

bin/s3_util: $(SRC)
	@echo "statically linking s3_util"
	CGO_ENABLED=0 GOOS=linux godep go build -a -installsuffix cgo -ldflags '-w -X main.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)' -o bin/s3_util cmd/s3_util/*.go

bin/sg_register: $(SRC)
	@echo "statically linking sg_register"
	CGO_ENABLED=0 GOOS=linux godep go build -a -installsuffix cgo -ldflags '-w -X main.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)' -o bin/sg_register cmd/sg_register/*.go

bootstrap:
	ginkgo bootstrap

bootstrap-test:
	@cd ${DIR} && ginkgo bootstrap && cd ~-

godeps:
	@echo "running godep"
	./go_get.sh

clean:
	@echo "running make clean"
	rm -f bin/awscli bin/ecr_login bin/ec2_tag
	docker images | grep -E '<none>' | awk '{print$$3}' | xargs docker rmi

distclean:
	@make clean
	@echo "running make distclean"
	rm -rf ./tmp ./certs
	docker rm awscli-build run-awscli
	docker rmi awscli-build $(REGISTRY)/awscli $(REGISTRY)/ecr_login $(REGISTRY)/ec2_tag $(REGISTRY)/sqs_util $(REGISTRY)/s3_util $(REGISTRY)/sg_register

interactive:
	@echo "running make interactive build"
	docker build -t awscli-build -f Dockerfile .
	docker run -it --rm --name awscli-build -v /var/run:/var/run -v $(MOUNT):/go/src/$(PROJ) --entrypoint=/bin/bash -i awscli-build

lint: $(SRC)
	@for pkg in $(LINTS); do echo "linting: $$pkg"; golint $$pkg; done

lint-check: $(SRC)
	@for pkg in $(LINTS); do \
		echo -n "linting: $$pkg: "; \
		echo "`golint $$pkg | wc -l | awk '{print$$NF}'` error(s)"; \
		[ $$(golint $$pkg | wc -l | awk '{print$$NF}') -le 0 ] && true || false; \
	done

run-awscli: config
	@echo "running bin/awscli"
	docker run -it --rm -i $(REGISTRY)/awscli

run-ecr_login: config
	@echo "running bin/ecr_login"
	docker run -it --rm -i $(REGISTRY)/ecr_login

run-ec2_tag: config
	@echo "running bin/ec2_tag"
	docker run -it --rm -i $(REGISTRY)/ec2_tag

run-sqs_util: config
	@echo "running bin/sqs_util"
	docker run -it --rm -i $(REGISTRY)/sqs_util

run-s3_util: config
	@echo "running bin/s3_util"
	docker run -it --rm -i $(REGISTRY)/s3_util

run-sg_register: config
	@echo "running bin/sg_register"
	docker run -it --rm -i $(REGISTRY)/sg_register

test:
	@echo "running test"
	docker run -it --rm -v /var/run:/var/run -v $(MOUNT):/go/src/$(PROJ) --entrypoint=/bin/sh -i awscli-build -c "godep restore && make test/units"

test/units: $(SRC)
	@echo "running test/units"
	godep go test $(TESTS) $(ARGS)

test-cover: $(SRC)
	@godep go test $(COVER) -coverprofile=coverage.out
	@godep go tool cover -html=coverage.out

vet: $(SRC)
	@for pkg in $(LINTS); do echo "vetting: $$pkg"; godep go vet $$pkg; done

.PHONY: clean test test/units run-bin/awscli run-bin/ecr_login run-bin/ec2_tag run-bin/sqs_util run-bin/s3_util run-bin/sg_register interactive bootstrap bootstrap-test lint lint-check test-cover godeps docker/push
