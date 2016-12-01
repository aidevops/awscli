FROM golang:1.7.3-wheezy

RUN apt-get update -y && \
  apt-get install -y vim && \
  curl -sL https://get.docker.com/builds/Linux/x86_64/docker-1.11.2.tgz -o /tmp/docker.tgz && \
  tar -xvzf /tmp/docker.tgz -C /usr/local && \
  ln -sf /usr/local/docker/docker /usr/bin/docker && \
  mkdir -m 700 /root/.ssh && \
  echo "Host github.com\n  StrictHostKeyChecking no" > /root/.ssh/config && \
  echo "[trusted]\ngroups = dialout" > /root/.hgrc && \
  go get github.com/tools/godep && \
  cp /go/bin/godep /usr/bin/godep && \
  go get github.com/onsi/ginkgo && \
  go get github.com/onsi/ginkgo/ginkgo && \
  cp /go/bin/ginkgo /usr/bin/ginkgo && \
  go get github.com/onsi/gomega && \
  go get -u github.com/golang/lint/golint && \
  cp /go/bin/golint /usr/bin/golint && \
  go get golang.org/x/tools/cmd/cover && \
  rm -rf /var/cache/apk/* /tmp/docker.tgz

WORKDIR /go/src/github.com/aidevops/awscli
