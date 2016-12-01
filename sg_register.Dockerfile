FROM scratch

ADD bin/sg_register /bin/sg_register
ADD certs /etc/ssl/certs

ENTRYPOINT [ "/bin/sg_register"]
CMD []
