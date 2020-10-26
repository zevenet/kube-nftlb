#!/usr/bin/env bash

set -e

for TEST_DIR in *; do
    if [ -d "$TEST_DIR" ] && [ "$TEST_DIR" != "template" ] && [ "$TEST_DIR" != "filters" ]; then
        kubectl delete --ignore-not-found -f "$TEST_DIR/input.yaml"
    fi
done

# No output means that it was already cleaned
