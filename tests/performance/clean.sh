#!/usr/bin/env bash

set -e

kubectl delete --ignore-not-found -f ./testdata/deployments/
kubectl delete --ignore-not-found -f ./testdata/services/

for KUBE_PATH in ./kubes/*
do
    kubectl delete --ignore-not-found -f "$KUBE_PATH"
done

# No output means that it was already cleaned
