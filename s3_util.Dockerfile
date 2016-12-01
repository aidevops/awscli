FROM scratch
#FROM alpine:edge

ADD bin/s3_util /bin/s3_util
ADD certs /etc/ssl/certs

#RUN apk add --update ca-certificates

ENTRYPOINT [ "/bin/s3_util"]
CMD []
