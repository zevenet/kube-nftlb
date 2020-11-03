#!/usr/bin/env bash

set -e

TIMEOUT=300s


##################
# Counting rules #
##################

function get_count_function {
    local KUBE_NAME=$1

    case $KUBE_NAME in
    kube-nftlb)
        echo -n "count_nftlb_rules"
        ;;
    kube-proxy)
        echo -n "count_iptables_rules"
        ;;
    *)
        echo "!! Unknown KUBE_NAME $KUBE_NAME, not implemented in get_count_function"
        return 1
        ;;
    esac
}

# kube-nftlb
function count_nftlb_rules {
    # nft list table ip nftlb   => List ruleset from nftlb table
    # | awk                     => Pipe ruleset as input text to awk
    # '/^\tchain/,/^\t}$/'      => Get "chain {...}" blocks (by nftables definition, they contain rules), 1 line == 1 rule
    # | sed                     => Pipe chain blocks as input text to sed
    # -e '/^\tchain/d'          => Delete "chain ..." lines
    # -e '/^\t}$/d'             => Delete "}" lines
    # -e '/^\t\ttype/d'         => Delete "type ..." lines (they're not rules)
    # -e '/^$/d'                => Delete empty lines
    # | wc                      => Pipe filtered (valid) rules as input text to wc
    # -l                        => Every rule is a line, so count every line

    nft list table ip nftlb | awk '/^\tchain/,/^\t}$/' | sed -e '/^\tchain/d' -e '/^\t}$/d' -e '/^\t\ttype/d' -e '/^$/d' | wc -l
}

# kube-proxy
function count_iptables_rules {
    # echo -n "$(tables)"   => Get legacy and actual iptables rulesets
    # | sed                 => Pipe ruleset as input text to sed
    # -e '/^#/d'            => Delete lines that start with "#"
    # -e '/^*/d'            => Delete lines that start with "*"
    # -e '/^:/d'            => Delete lines that start with ":"
    # -e '/^COMMIT/d'       => Delete "COMMIT" lines
    # -e '/^$/d'            => Delete empty lines
    # | wc                  => Pipe filtered (valid) rules as input text to wc
    # -l                    => Every rule is a line, so count every line

    echo -n "$(iptables-save)$(iptables-legacy-save)" | sed -e '/^#/d' -e '/^*/d' -e '/^:/d' -e '/^COMMIT/d' -e '/^$/d' | wc -l

    # Sleep for 10ms not to block iptables, nftables doesn't have this problem
    sleep 0.01
}


##########
# Timers #
##########

function get_timer_function {
    local TIMER_TYPE=$1

    case $TIMER_TYPE in
    create)
        echo -n "timer_rule_count_increase"
        ;;
    delete)
        echo -n "timer_rule_count_decrease"
        ;;
    *)
        echo "!! Unknown TIMER_TYPE $TIMER_TYPE, not implemented in get_timer_function"
        return 1
        ;;
    esac
}

function timer_rule_count_increase {
    local COUNT_FUNCTION=$1
    local EXPECTED_RULE_COUNT=$2

    local RULE_COUNT=$($COUNT_FUNCTION)
    while [ $RULE_COUNT -lt $EXPECTED_RULE_COUNT ]
    do
        RULE_COUNT=$($COUNT_FUNCTION)
    done
}

function timer_rule_count_decrease {
    local COUNT_FUNCTION=$1
    local EXPECTED_RULE_COUNT=$2

    local RULE_COUNT=$($COUNT_FUNCTION)
    while [ $RULE_COUNT -gt $EXPECTED_RULE_COUNT ]
    do
        RULE_COUNT=$($COUNT_FUNCTION)
    done
}

function timer_show {
    local KUBE_NAME=$1
    local RESOURCE_NAME=$2
    local COUNT_TYPE=$3
    local TIMER_TYPE=$(echo "$COUNT_TYPE" | sed 's/-.*$//g')

    # Get expected rule count for this resource
    local EXPECTED_RULE_COUNT=$(grep "$COUNT_TYPE" "./testdata/expected-rule-count/$KUBE_NAME/$RESOURCE_NAME.txt" | sed "s/$COUNT_TYPE: //")

    # What functions are going to be executed for counting and timing?
    local COUNT_FUNCTION=$(get_count_function "$KUBE_NAME")
    local TIMER_FUNCTION=$(get_timer_function "$TIMER_TYPE")

    # Timer explanation:
    # date +%s      =>  Returns seconds without decimals.
    # date +%N      =>  Returns nanoseconds in this actual second.
    # date +%s.%N   =>  Returns "seconds.nanoseconds", which is nice but we can't do floating point math in bash easily.
    # date +%s%N    =>  Returns "secondsnanoseconds". Those seconds can be simplified as nanoseconds without the dot (as
    #                   if we were multiplying it by 10^9) and then the remaining nanoseconds is added to the result.
    #                   At last, if we divide that result by 10^6, we get milliseconds without decimals. 10^9/10^6 = 10^3,
    #                   so if we keep the first 3 digits from nanoseconds we get milliseconds.
    # date +%s%3N   =>  Returns milliseconds without decimals.
    local TIME_START=$(date +%s%3N)
    "$TIMER_FUNCTION" "$COUNT_FUNCTION" "$EXPECTED_RULE_COUNT"
    local TIME_END=$(date +%s%3N)

    # Show results
    local TIME_RESULT=$(( TIME_END - TIME_START ))
    echo "$COUNT_TYPE: $TIME_RESULT ms ($EXPECTED_RULE_COUNT rules counted)"
}


###############
# Deployments #
###############

function create_deployment {
    local DEPLOYMENT_PATH=$1
    local RESOURCE_NAME=$2

    kubectl apply -f "$DEPLOYMENT_PATH" --timeout="$TIMEOUT"
    while [ $(kubectl get pods -A | grep "$RESOURCE_NAME" | grep -v Running | wc -l) -gt 0 ]; do
        sleep 1
    done
    sleep 5
}

function delete_deployment {
    local DEPLOYMENT_PATH=$1
    local RESOURCE_NAME=$2

    kubectl delete -f "$DEPLOYMENT_PATH" --timeout="$TIMEOUT"
    kubectl wait --for=delete pods -l app="$RESOURCE_NAME" --timeout="$TIMEOUT"
}


############
# Services #
############

function create_service {
    local SERVICE_PATH=$1

    kubectl apply -f "$SERVICE_PATH" --timeout="$TIMEOUT"
}

function delete_service {
    local SERVICE_PATH=$1

    kubectl delete -f "$SERVICE_PATH" --timeout="$TIMEOUT"
}


##############################
# kube: create, test, delete #
##############################

function create_kube {
    local KUBE_NAME=$1
    local KUBE_PATH=$2

    echo "Starting $KUBE_NAME..."
    kubectl apply -f "$KUBE_PATH" --timeout="$TIMEOUT"
    sleep 3
    kubectl wait --namespace=kube-system --for=condition=Ready pods -l app="$KUBE_NAME" --timeout="$TIMEOUT"
}

function delete_kube {
    local KUBE_NAME=$1
    local KUBE_PATH=$2

    echo "Deleting $KUBE_NAME..."
    kubectl delete -f "$KUBE_PATH" --timeout="$TIMEOUT"
    sleep 3
    kubectl wait --namespace=kube-system --for=delete pods -l app="$KUBE_NAME" --timeout="$TIMEOUT"
}

# Test every deployment given a kube name and save how much time is spent in creating and deleting its services.
# Also, print how many rules are after creating and deleting services.
function test_kube {
    local KUBE_NAME=$1

    for DEPLOYMENT_PATH in ./testdata/deployments/*
    do
        # Useful values as parameters
        local RESOURCE_NAME=$(echo "$DEPLOYMENT_PATH" | sed -e 's/^.*testdata\/deployments\///' -e 's/[.]yaml$//')
        local SERVICE_PATH="./testdata/services/$RESOURCE_NAME.yaml"

        # Create deployment and wait for it to be created
        echo -e "\t$RESOURCE_NAME {"
        create_deployment "$DEPLOYMENT_PATH" "$RESOURCE_NAME"

        # Create service and time it until it reaches the expected (increased) rule count
        create_service "$SERVICE_PATH"
        timer_show "$KUBE_NAME" "$RESOURCE_NAME" "create-service"

        # Delete service and time it until it reaches the expected (decreased) rule count
        delete_service "$SERVICE_PATH"
        timer_show "$KUBE_NAME" "$RESOURCE_NAME" "delete-service"

        # Delete deployment and wait for it to be deleted, also time it until it reaches the expected (decreased) rule count
        delete_deployment "$DEPLOYMENT_PATH" "$RESOURCE_NAME"
        echo -e "\t}"
    done
}


####################
# Main script loop #
####################

# For every kube dir, get its path
for KUBE_PATH in ./kubes/*
do
    # Main parameter (extract name from kube path)
    KUBE_NAME=$(echo "$KUBE_PATH" | sed -e 's/^.*kubes\///' -e 's/[.]yaml$//')

    # Create kube daemonset and apply its configuration, and wait for it to be created
    echo "$KUBE_NAME {"
    create_kube "$KUBE_NAME" "$KUBE_PATH"

    # Hardcoded grace time
    sleep 30

    # Test deployments and services using this kube
    test_kube "$KUBE_NAME"

    # Delete kube daemonset and its configuration, and wait for it to be deleted
    delete_kube "$KUBE_NAME" "$KUBE_PATH"
    echo "}"
done
