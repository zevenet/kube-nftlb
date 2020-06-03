#!/bin/sh

DOCKER_PATH="./internal/docker/kube-nftlb"

# Optionally, use a nftlb devel package
if [ ! -z $1 ]; then
	cp $1 $DOCKER_PATH/nftlb.deb
else
	# use empty file to avoid docker COPY directive failure
	touch $DOCKER_PATH/nftlb.deb
fi


# It lets you interact with nftlb.
# The key must be propagated to the client and the daemon.
key="12345"

# Before running this script, you MUST BE ROOT and you need to have the following tools installed:
#   - Docker
#   - Docker-machine
#   - Minikube (Kubernetes local cluster, virtualized)
#   - Golang
#   - client-go libs

# STEP 1:
#   Compile cmd/app/main.go.
#   The binary will be called "app".
GOOS=linux go build -o $DOCKER_PATH/app ./cmd/app

# YOU MUST DO THIS STEP IF YOU DIDN'T STARTED MINIKUBE WITH 'minikube start --vm-driver=none':
#   Uncomment the line below if you apply.
#eval $(minikube docker-env)

# STEP 2:
#   The client container will be created using its Dockerfile.
#   It will be made for Docker, not for Minikube (this will come later).
docker build -t kube-nftlb $DOCKER_PATH --build-arg KEY=$key

# STEP 3:
#   Clean residual files.
rm -f $DOCKER_PATH/app
rm $DOCKER_PATH/nftlb.deb

# STEP 4:
#   Every nftables rule should be flushed before creating kube-nftlb.
#   Do it at your own risk (you will probably have trouble with Docker).
#nft flush ruleset
