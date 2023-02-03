FROM amazon/aws-cli:amd64

LABEL maintainer="OP (ompragash) <ompragash@proton.me>"

ADD release/linux/amd64/drone-jacoco-s3-upload-publish /bin/

ENTRYPOINT ["/usr/bin/bash", "-l", "-c", "/bin/drone-jacoco-s3-upload-publish"]
