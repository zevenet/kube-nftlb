#!/bin/sh

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
GOOS=linux go build -o ./internal/docker/client/app ./cmd/app

# YOU MUST DO THIS STEP IF YOU DIDN'T STARTED MINIKUBE WITH 'minikube start --vm-driver=none':
#   Uncomment the line below if you apply.
#eval $(minikube docker-env)

# STEP 2:
#   The client container will be created using its Dockerfile.
#   It will be made for Docker, not for Minikube (this will come later).
docker build -t client internal/docker/client --build-arg KEY=$key

# STEP 3:
#   The daemon container will be created using its Dockerfile.
#   It will be made for Docker, not for Minikube (this will come later).
docker build -t daemon internal/docker/daemon --build-arg KEY=$key

# STEP 4:
#   Clean residual files.
rm -f docker/client/app
