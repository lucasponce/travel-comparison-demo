#!/bin/bash

DOCKER_NAME=lucasponce/travels
DOCKER_VERSION=dev
DOCKER_TAG=${DOCKER_NAME}:${DOCKER_VERSION}

go build -o docker/travels travels.go

docker build -t ${DOCKER_TAG} docker
