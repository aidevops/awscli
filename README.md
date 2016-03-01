awscli
------

** THIS IS AN EXPERIMENT IN CONTAINER SIZE REDUCTION, USE AT YOUR OWN RISK. **
Comparing image sizes and efficiency gained.....:
------------------------------------------------

The first attempt
=================

At delivering a containerized tag instance was in the range of ~300 - ~900 MB with debian/ubuntu, centos, fedora images.


My second attempt
=================

Looked something like this at around ~113 - ~150 MB with alpine. It included a base layer and a shell script tag overlay layer.



- The first part:

``` Dockerfile
# awscli - minimal container to run aws cli tools
#
# VERSION               0.0.2
FROM alpine:latest
MAINTAINER John Torres <john.torres@pearson.com>

RUN \
  apk add --update python git curl unzip groff && \
  curl -s "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" -o "awscli-bundle.zip" && \
  unzip awscli-bundle.zip -d /tmp/ && \
  /tmp/awscli-bundle/install -i /usr/local/aws -b /usr/local/bin/aws && \
  apk del unzip curl && \
  rm -rf /tmp/awscli-bundle* /var/cache/apk/*


WORKDIR /root

ENTRYPOINT ["/bin/sh"]

CMD []

```

- The second part (overlay):


``` Dockerfile
# awstag - submit a tag and exit
#
# VERSION               0.0.2

FROM aidevops/awscli:0.0.2
MAINTAINER John Torres <john.torres@pearson.com>

# install tagging script
ADD bin/aws_tag.sh /etc/aws_tag.sh

RUN chmod 500 /etc/aws_tag.sh

ENTRYPOINT ["/etc/aws_tag.sh"]

WORKDIR /root

```

- Where `aws_tag.sh` is:

```
#!/bin/sh
#
aws ec2 --region=${region} create-tags \
  --resource ${instance_id} \
  --tags \
    "Key=build,Value=$name-$version.$dns_domain-$os-$os_version-$instance_size@$cluster_size" \
    "Key=environment,Value=$customer-$consortium-$environment" \
    "Key=region,Value=$region" \
    "Key=consul_url,Value=$consul_url" \
    "Key=consul_dc,Value=$consul_dc" \
    "Key=role,Value=$role" \
    "Key=node_name,Value=$node_name" \
    "Key=Name,Value=$name" \
    "Key=instance_id,Value=$instance_id" \
    "Key=etcd_discovery,Value=$ETCD_DISCOVERY"
```


My third Attempt
================

Is to use the aws-go-sdk and package static binaries. The above `aws_tag.sh` will
be augmented slightly to reflect the go binary flags.

See below for more detail.


Quick Start:
-----------

- Setup environment

  `mkdir -p awscli/src/github.com/johnt337`

  `git clone github.com/johnt337/awscli awscli/src/github.com/johnt337/awscli`

  `cd github.com/johnt337/awscli awscli/src/github.com/johnt337/awscli && ./dev_setup.sh`



- Build the build environment followed by the entire suite

  `cd github.com/johnt337/awscli awscli/src/github.com/johnt337/awscli`

  `make build`

- Build just the `ecr_login` util

  `cd github.com/johnt337/awscli awscli/src/github.com/johnt337/awscli`

  `make build-ecr_login`

- Build just the `ec2_tag` util

  `cd github.com/johnt337/awscli awscli/src/github.com/johnt337/awscli`

  `make build-ec2_tag`


Running:
-------

- Run tag similar to aws ecr get-login --region <region> --registry-ids <id1,id2,id3> 

  `eval $(docker run --rm -it johnt337/ecr_login -account=$AWS_REGISTRY_ID)`

- Run tag with the bundled docker

  `docker run --rm -it -v $HOME/.docker:/root/.docker -v /var/run/docker.sock:/var/run/docker.sock -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY johnt337/ecr_login -account=$AWS_REGISTRY_ID -login`
