FROM alpine

ADD bin/ecr_login /bin/ecr_login
ADD certs /etc/ssl/certs

RUN apk add --no-cache curl && \
    curl -sL https://get.docker.com/builds/Linux/x86_64/docker-1.11.2.tgz -o /tmp/docker.tgz && \
    tar -xvzf /tmp/docker.tgz -C /usr/local && \
    ln -sf /usr/local/docker/docker /usr/bin/docker && \
    chmod +x /usr/bin/docker && \
    apk del --no-cache curl && \
    rm -f /tmp/docker.tgz

ENTRYPOINT [ "/bin/ecr_login"]
CMD []
