#!/usr/bin/env bash

set -e

if [ -z "$1" ]
then
    echo -e '/!\ Missing first argument: a kube name'
    exit 1
elif [ -z "$2" ]
then
    echo -e '/!\ Missing second argument: a kube name'
    exit 1
fi

KUBE_NAME_1=$1
KUBE_NAME_2=$2
POOL_BASE_PATH="pool"
CHARTS_BASE_PATH="charts"
Y_AXES=("time" "rules")
COUNT_TYPES=("create-service" "delete-service")


#############
# Functions #
#############

function get_path {
    local BASE_PATH=$1
    local Y_AXIS=$2
    local X_AXIS=$3

    # This follows the tree structure in pool/ and charts/
    echo -n "$BASE_PATH/$Y_AXIS/$X_AXIS"
}

function get_script_path {
    local CHARTS_PATH=$1

    echo -n "$CHARTS_PATH/chart.gp"
}

function get_output_path {
    local CHARTS_PATH=$1
    local RESOURCE_POOL_NAME=$2

    echo -n "$CHARTS_PATH/$RESOURCE_POOL_NAME.png"
}

function get_data_path {
    local POOL_PATH=$1
    local KUBE_NAME=$2
    local RESOURCE_POOL_NAME=$3

    echo -n "$POOL_PATH/$KUBE_NAME-$RESOURCE_POOL_NAME.txt"
}

function get_chart_title {
    local X_AXIS_DESCRIPTION=$1
    local TITLE_DESCRIPTION=$2

    case $X_AXIS_DESCRIPTION in
    create-service)
        X_AXIS_DESCRIPTION="Creating a Service"
    ;;
    delete-service)
        X_AXIS_DESCRIPTION="Deleting a Service"
    ;;
    esac

    # Better titles for charts
    echo -n "$X_AXIS_DESCRIPTION $TITLE_DESCRIPTION"
}

function generate_chart {
    local RESOURCE_POOL_NAME=$1
    local Y_AXIS=$2
    local X_AXIS=$3
    local X_AXIS_DESCRIPTION=$4
    local TITLE_DESCRIPTION=$5

    # Get pool and charts full path
    local POOL_PATH=$(get_path "$POOL_BASE_PATH" "$Y_AXIS" "$X_AXIS")
    local CHARTS_PATH=$(get_path "$CHARTS_BASE_PATH" "$Y_AXIS" "$X_AXIS")

    # Get data needed for the script
    local SCRIPT_FILE=$(get_script_path "$CHARTS_PATH")
    local CHART_TITLE=$(get_chart_title "$X_AXIS_DESCRIPTION" "$TITLE_DESCRIPTION")
    local OUTPUT_FILE=$(get_output_path "$CHARTS_PATH" "$RESOURCE_POOL_NAME")
    local DATA_PATH_1=$(get_data_path "$POOL_PATH" "$KUBE_NAME_1" "$RESOURCE_POOL_NAME")
    local DATA_PATH_2=$(get_data_path "$POOL_PATH" "$KUBE_NAME_2" "$RESOURCE_POOL_NAME")

    # Generate chart
    gnuplot -c "$SCRIPT_FILE" "$CHART_TITLE" "$OUTPUT_FILE" "$DATA_PATH_1" "$DATA_PATH_2" "$KUBE_NAME_1" "$KUBE_NAME_2"
}


#######################
# Generate all charts #
#######################

for DEPLOYMENT_PATH in ./testdata/deployments/*
do
    RESOURCE_NAME=$(echo "$DEPLOYMENT_PATH" | sed -e 's/^.*testdata\/deployments\///' -e 's/[.]yaml$//')
    NUMBER_ENDPOINTS=$(echo "$RESOURCE_NAME" | sed -e 's/.*-//g' -e 's/^0*//g')

    for Y_AXIS in ${Y_AXES[*]}
    do
        generate_chart "$RESOURCE_NAME" "$Y_AXIS" "by-count-type" "$NUMBER_ENDPOINTS" "endpoints"

        for COUNT_TYPE in ${COUNT_TYPES[*]}
        do
            generate_chart "$COUNT_TYPE" "$Y_AXIS" "by-endpoints-number" "$COUNT_TYPE" "for N endpoints"
        done
    done
done
