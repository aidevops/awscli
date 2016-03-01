FROM scratch

ADD bin/ecr_login /bin/ecr_login
ADD certs /etc/ssl/certs

ENTRYPOINT [ "/bin/ecr_login"]
CMD []