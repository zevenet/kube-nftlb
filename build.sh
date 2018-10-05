#!/bin/bash

# Before running this script, you MUST BE ROOT and you need to have the following tools installed:
#   - Docker
#   - Docker-machine
#   - Minikube (Kubernetes local cluster, virtualized)
#   - Golang
#   - client-go libs

# STEP 1: 
#   Compile go/main.go.
#   The binary will be called "app".
GOOS=linux go build -o ./app ./go

# YOU MUST DO THIS STEP IF YOU DIDN'T STARTED MINIKUBE WITH 'minikube start --vm-driver=none':
#   Uncomment the line below if you apply.
#eval $(minikube docker-env)

# STEP 2:
#   A container will be created using the attached Dockerfile.
#   It will be made for Docker, not for Minikube (this will come later).
docker build -t nftlb .
