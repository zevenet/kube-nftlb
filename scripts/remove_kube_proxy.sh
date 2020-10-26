#!/usr/bin/env bash

kubectl delete --ignore-not-found -f tests/performance/kubes/kube-proxy.yaml
