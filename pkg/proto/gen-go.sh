#!/bin/bash

# Developer in charge of this process (not part of build)
# see README.md for what is required to have installed in your machine (docker)
set -eo pipefail
set -x

# Path to proto files

PROTO_FILE=./pkg/proto/user.proto

DOCKER_IMAGE=devlube/protogen:3.0.0


# Clean up previously generated
[ -e "$PROTO_FILE.pb.*" ] && rm $PROTO_FILE.pb.*


# Generates grpc server code
docker run -v $PWD:/tmp -w /tmp $DOCKER_IMAGE protoc \
    -I=/go/src:. \
    -I/thirdparty/googleapis \
    --lint_out=. \
    --gogo_out=plugins=grpc,paths=source_relative:. \
    $PROTO_FILE

