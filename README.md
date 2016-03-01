awscli
------

** THIS IS AN EXPERIMENT IN CONTAINER SIZE REDUCTION, USE AT YOUR OWN RISK.


Quick Start:
-----------

- Setup environment

  `mkdir -p awscli/src/github.com/johnt337`

  `git clone github.com/johnt337/awscli awscli/src/github.com/johnt337/awscli`

  `cd github.com/johnt337/awscli awscli/src/github.com/johnt337/awscli && ./dev_setup.sh`



- Build the build environment followed by the entire suite

  `cd github.com/johnt337/awscli awscli/src/github.com/johnt337/awscli`

  `make`

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
