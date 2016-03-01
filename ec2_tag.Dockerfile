FROM scratch

ADD bin/ec2_tag /bin/ec2_tag
ADD certs /etc/ssl/certs

ENTRYPOINT [ "/bin/ec2_tag"]
CMD []