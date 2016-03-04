FROM scratch

ADD bin/sqs_util /bin/sqs_util
ADD certs /etc/ssl/certs

ENTRYPOINT [ "/bin/sqs_util"]
CMD []
