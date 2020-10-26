#!/usr/bin/env bash
# shellcheck disable=SC1091

# Root directory is 2 levels up
source ../../.env

echo "Obtaining expected data for test $1..."

# Apply YAML resources
kubectl apply -f "$1/input.yaml" >/dev/null
sleep 15

for FILENAME_JSON in "$1"/*.json; do
    # Get farm name to read the expected JSON
    FARM_NAME=$(echo "$FILENAME_JSON" | sed -f filters/get-farm-name-from-filename.sed)

    # Get expected JSON and save that as a file
    JSON=$(curl -s -H "Key: $NFTLB_KEY" "$NFTLB_PROTOCOL://$NFTLB_HOST:$NFTLB_PORT/farms/$FARM_NAME" | sed -f filters/replace-farm-values.sed | jq --indent 4 -S .)
    echo -n "$JSON" > "$FILENAME_JSON"
done

# Get expected nft ruleset and save that as a file
NFT_RULESET=$(nft list table ip nftlb | awk -f filters/select-chains-nft-ruleset.awk | sed -f filters/clean-chains-nft-ruleset.sed | sort)
echo -n "$NFT_RULESET" > "$1/ruleset.nft"

# Delete YAML resources
kubectl delete -f "$1/input.yaml" >/dev/null
