#!/bin/sh

# /!\ Before running this script, you MUST BE ROOT and Docker must be installed

set -e

DOCKER_PATH="./docker"

# Optionally, use a nftlb devel package
if [ ! -z $1 ]; then
	cp $1 $DOCKER_PATH/nftlb.deb
else
	# Use empty file to avoid docker COPY directive failure
	touch $DOCKER_PATH/nftlb.deb
fi

# Uncomment the line below if you didn't start Minikube with 'minikube start --vm-driver=none'.
#eval $(minikube docker-env)

# The container image will be built using its Dockerfile. Minikube will use this image later to make the container.
docker build -t kube-nftlb $DOCKER_PATH

# Clean residual files and intermediate containers.
rm $DOCKER_PATH/nftlb.deb
docker image prune -f --filter label=stage=intermediate

# Every nftlb rule should be flushed before creating kube-nftlb. Do it at your own risk.
#nft flush table nftlb
