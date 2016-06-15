FROM scratch

ADD bin/s3_util /bin/s3_util
ADD certs /etc/ssl/certs

ENTRYPOINT [ "/bin/s3_util"]
CMD []
