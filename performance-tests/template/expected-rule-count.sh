#!/usr/bin/env bash

set -e

TIMEOUT=180s
SLEEP_TIME=30

function count_rules {
    # You are expected to put the logic here to get the raw count of rules
}

function print_count {
    # This prints the result with the correct format, you are not supposed to change this
    echo "$1: $(count_rules)"
}


########
# Main #
########

# Params
KUBE_PATH=$1
KUBE_NAME=$(echo "$KUBE_PATH" | sed -e 's/^[.]\/kubes\///' -e 's/[.]yaml$//')

# Create kube
kubectl apply -f "$KUBE_PATH" --timeout="$TIMEOUT"
kubectl wait --namespace=kube-system --for=condition=Ready pods -l app="$KUBE_NAME" --timeout="$TIMEOUT"
print_count "create-kube"

# Test every deployment and its service
for DEPLOYMENT_PATH in ./testdata/deployments/*
do
    DEPLOYMENT_NAME=$(echo "$DEPLOYMENT_PATH" | sed -e 's/^[.]\/testdata\/deployments\///' -e 's/[.]yaml$//')
    SERVICE_PATH=$(echo "$DEPLOYMENT_PATH" | sed -e 's/deployments/services/')

    echo "=> Deployment $DEPLOYMENT_NAME"

    # Create deployment
    kubectl apply -f "$DEPLOYMENT_PATH" --timeout="$TIMEOUT"
    kubectl wait -f "$DEPLOYMENT_PATH" --for=condition=Available --timeout="$TIMEOUT"
    sleep "$SLEEP_TIME"
    print_count "create-deployment"

    # Create service
    kubectl apply -f "$SERVICE_PATH" --timeout="$TIMEOUT"
    sleep "$SLEEP_TIME"
    print_count "create-service"

    # Delete service
    kubectl delete -f "$SERVICE_PATH" --timeout="$TIMEOUT"
    sleep "$SLEEP_TIME"
    print_count "delete-service"

    # Delete deployment
    kubectl delete -f "$DEPLOYMENT_PATH" --timeout="$TIMEOUT"
    kubectl wait --for=delete pods -l app="$DEPLOYMENT_NAME" --timeout="$TIMEOUT"
    sleep "$SLEEP_TIME"
    print_count "delete-deployment"

    echo
done

# Delete kube
kubectl delete -f "$KUBE_PATH" --timeout="$TIMEOUT"
kubectl wait --namespace=kube-system --for=delete pods -l app="$KUBE_NAME" --timeout="$TIMEOUT"
