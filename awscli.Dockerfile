FROM scratch

ADD bin/awscli        /bin/awscli

#ADD config /etc/awscli
ADD certs /etc/ssl/certs
ADD tmp /tmp

ENTRYPOINT [ "/bin/awscli"]
CMD []
