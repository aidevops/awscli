FROM alpine

ADD https://get.docker.com/builds/Linux/x86_64/docker-1.8.3 /usr/bin/docker

ADD bin/ecr_login /bin/ecr_login
ADD certs /etc/ssl/certs

RUN chmod +x /usr/bin/docker

ENTRYPOINT [ "/bin/ecr_login"]
CMD []