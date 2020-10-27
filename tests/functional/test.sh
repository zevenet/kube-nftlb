#!/usr/bin/env bash
# shellcheck disable=SC1091

# Root directory is 2 levels up
source ../../.env

# Read args, this lets us test individual cases
if [ $# -eq 0 ]; then
    DIRS=*
else
    DIRS="$@"
fi

# Read every file or directory
for TEST_DIR in $DIRS; do
    # Skip this iteration if it isn't a directory or it's blacklisted
    if [ ! -d "$TEST_DIR" ] || [ "$TEST_DIR" = "template" ] || [ "$TEST_DIR" = "filters" ]; then
        continue
    fi

    echo "→ Running test $TEST_DIR..."
    TEST_PASSED="true"

    # Apply YAML resources
    kubectl apply -f "$TEST_DIR/input.yaml" >/dev/null
    sleep 45

    # Read filenames from .json files
    for FILENAME_JSON in "$TEST_DIR"/*.json; do
        # Get farm name to read the expected JSON
        FARM_NAME=$(echo "$FILENAME_JSON" | sed -f filters/get-farm-name-from-filename.sed)

        # Get actual and expected JSON
        JSON=$(curl -s -H "Key: $NFTLB_KEY" "$NFTLB_PROTOCOL://$NFTLB_HOST:$NFTLB_PORT/farms/$FARM_NAME" | sed -f filters/replace-farm-values.sed | jq --indent 4 -S .)
        EXPECTED_JSON=$(cat "$FILENAME_JSON" | sed -f filters/replace-farm-values.sed | jq --indent 4 -S .)

        # Compare both JSON strings
        if [ "$JSON" != "$EXPECTED_JSON" ]; then
            echo "🚨 Error in farm $FARM_NAME, the JSON doesn't match the expected result"
            diff --color -u <(echo "$EXPECTED_JSON") <(echo "$JSON")
            TEST_PASSED="false"
            echo # Separation line if there's more than 1 error
        fi
    done

    # Get actual and expected nft rulesets
    NFT_RULESET=$(echo -n "$(nft list table ip nftlb)$(nft list table ip netdev 2>/dev/null)" | awk -f filters/select-chains-nft-ruleset.awk | sed -f filters/clean-chains-nft-ruleset.sed | sort)
    EXPECTED_NFT_RULESET=$(cat "$TEST_DIR/ruleset.nft" | awk -f filters/select-chains-nft-ruleset.awk | sed -f filters/clean-chains-nft-ruleset.sed | sort)

    # Compare both rulesets
    if [ "$NFT_RULESET" != "$EXPECTED_NFT_RULESET" ]; then
        echo "🚨 nft ruleset after applying input.yaml doesn't match the expected result"
        diff --color -u <(echo "$EXPECTED_NFT_RULESET") <(echo "$NFT_RULESET")
        TEST_PASSED="false"
        echo # Empty line
    fi

    # Delete YAML resources
    kubectl delete -f "$TEST_DIR/input.yaml" >/dev/null
    sleep 15

    # Check ruleset against the clean one
    NFT_RULESET=$(echo -n "$(nft list table ip nftlb)$(nft list table ip netdev 2>/dev/null)" | awk -f filters/select-chains-nft-ruleset.awk | sed -f filters/clean-chains-nft-ruleset.sed | sort)
    EXPECTED_NFT_RULESET=$(cat clean-ruleset.nft | awk -f filters/select-chains-nft-ruleset.awk | sed -f filters/clean-chains-nft-ruleset.sed | sort)

    if [ "$NFT_RULESET" != "$EXPECTED_NFT_RULESET" ]; then
        echo "🚨 nft ruleset after deleting input.yaml doesn't match the expected result"
        diff --color -u <(echo "$EXPECTED_NFT_RULESET") <(echo "$NFT_RULESET")
        TEST_PASSED="false"
        echo # Empty line
    fi

    if [ "$TEST_PASSED" = "true" ]; then
        echo "✅ TEST PASSED"
    else
        echo "❌ TEST FAILED"
    fi
    echo # Empty line
done
