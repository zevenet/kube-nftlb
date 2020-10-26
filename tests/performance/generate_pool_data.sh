#!/usr/bin/env bash

set -e

if [ -z "$1" ]
then
    echo -e '/!\ Missing filtered-results.txt'
    exit 1
fi

FILE_PATH=$1
POOL_BASE_PATH="pool"
POOL_TIME_PATH="$POOL_BASE_PATH/time"
POOL_RULES_PATH="$POOL_BASE_PATH/rules"
COUNT_TYPES=("create-service" "delete-service")


################
# y-axis: time #
################
for KUBE_PATH in ./kubes/*
do
    KUBE_NAME=$(echo "$KUBE_PATH" | sed -e 's/^.*kubes\///' -e 's/[.]yaml$//')

    #########################
    # x-axis: by-count-type #
    #########################
    for DEPLOYMENT_PATH in ./testdata/deployments/*
    do
        RESOURCE_NAME=$(echo "$DEPLOYMENT_PATH" | sed -e 's/^.*testdata\/deployments\///' -e 's/[.]yaml$//')
        COLUMN_PATHS=()

        for COUNT_TYPE in ${COUNT_TYPES[*]}
        do
            FULL_COLUMN_PATH="$POOL_TIME_PATH/columns/$KUBE_NAME-$RESOURCE_NAME-$COUNT_TYPE.txt"
            COLUMN_PATHS+=("$FULL_COLUMN_PATH")

            # time/columns are generated here
            awk "/^$KUBE_NAME {$/,/^}$/" "$FILE_PATH" | awk "/^\\t$RESOURCE_NAME {$/{flag=1;next}/^\\t}$/{flag=0}flag" | grep "$COUNT_TYPE" | sed -e "s/^$COUNT_TYPE: //g" -e 's/ ms .*$//g' > "$FULL_COLUMN_PATH"
        done

        paste -d' ' ${COLUMN_PATHS[*]} > "$POOL_TIME_PATH/by-count-type/$KUBE_NAME-$RESOURCE_NAME.txt"
    done

    ###############################
    # x-axis: by-endpoints-number #
    ###############################
    for COUNT_TYPE in ${COUNT_TYPES[*]}
    do
        COLUMN_PATHS=()

        for DEPLOYMENT_PATH in ./testdata/deployments/*
        do
            RESOURCE_NAME=$(echo "$DEPLOYMENT_PATH" | sed -e 's/^.*testdata\/deployments\///' -e 's/[.]yaml$//')
            FULL_COLUMN_PATH="$POOL_TIME_PATH/columns/$KUBE_NAME-$RESOURCE_NAME-$COUNT_TYPE.txt"
            COLUMN_PATHS+=("$FULL_COLUMN_PATH")
        done

        paste -d' ' ${COLUMN_PATHS[*]} > "$POOL_TIME_PATH/by-endpoints-number/$KUBE_NAME-$COUNT_TYPE.txt"
    done
done


#################
# y-axis: rules #
#################
for KUBE_EXPECTED_RULE_COUNT_PATH in ./testdata/expected-rule-count/*
do
    KUBE_NAME=$(echo "$KUBE_EXPECTED_RULE_COUNT_PATH" | sed -e 's/^.*testdata\/expected-rule-count\///' -e 's/\/.*$//')

    #########################
    # x-axis: by-count-type #
    #########################
    for KUBE_RESOURCE_COUNT_FILE in $KUBE_EXPECTED_RULE_COUNT_PATH/*
    do
        RESOURCE_NAME=$(echo "$KUBE_RESOURCE_COUNT_FILE" | sed -e 's/^.*\///' -e 's/[.]txt$//')
        COLUMN_PATHS=()

        for COUNT_TYPE in ${COUNT_TYPES[*]}
        do
            FULL_COLUMN_PATH="$POOL_RULES_PATH/columns/$KUBE_NAME-$RESOURCE_NAME-$COUNT_TYPE.txt"
            COLUMN_PATHS+=("$FULL_COLUMN_PATH")

            # rules/columns are generated here
            grep "$COUNT_TYPE" "$KUBE_RESOURCE_COUNT_FILE" | sed "s/^$COUNT_TYPE: //" > "$FULL_COLUMN_PATH"
        done

        paste -d'\n' ${COLUMN_PATHS[*]} > "$POOL_RULES_PATH/by-count-type/$KUBE_NAME-$RESOURCE_NAME.txt"
    done

    ###############################
    # x-axis: by-endpoints-number #
    ###############################
    for COUNT_TYPE in ${COUNT_TYPES[*]}
    do
        COLUMN_PATHS=()

        for KUBE_RESOURCE_COUNT_FILE in $KUBE_EXPECTED_RULE_COUNT_PATH/*
        do
            RESOURCE_NAME=$(echo "$KUBE_RESOURCE_COUNT_FILE" | sed -e 's/^.*\///' -e 's/[.]txt$//')
            FULL_COLUMN_PATH="$POOL_RULES_PATH/columns/$KUBE_NAME-$RESOURCE_NAME-$COUNT_TYPE.txt"
            COLUMN_PATHS+=("$FULL_COLUMN_PATH")
        done

        paste -d'\n' ${COLUMN_PATHS[*]} > "$POOL_RULES_PATH/by-endpoints-number/$KUBE_NAME-$COUNT_TYPE.txt"
    done
done
