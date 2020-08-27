#!/bin/bash

# Test ExternalIPs

NFTLB_KEY=$(grep 'NFTLB_KEY' ../../.env | sed 's/NFTLB_KEY=//')

# Create Service
kubectl apply -f . &>/dev/null
sleep 10
curl --silent -H "Key: $NFTLB_KEY" http://localhost:5555/farms > configCreation.nft

# Delete Service
kubectl delete -f . &>/dev/null
sleep 10
curl --silent -H "Key: $NFTLB_KEY" http://localhost:5555/farms > configDelete.nft
# Apply format
# We substitute all the value of the virtual-addr fields except for the externalIPs that we have defined in the service.
# If we change these values, we must manually modify the range of ips (192.168.10.[89-90-91])

sed -i '/\"192[.]168[.]10[.]\(89\|90\|91\)\"/!s/\("virtual-addr": \)\(.*\)/\1"IP"/' configCreation.nft
sed -i 's/\("ip-addr": \)\(.*\)/\1"IP"/' configCreation.nft
sed -i -r 'N;s/.*(\n(\s*)"ip-addr":)/\2"name": "NAME",\1/;P;D' configCreation.nft
sed -i 's/\("virtual-addr": \)\(.*\)/\1"IP"/' configDelete.nft
sed -i 's/\("ip-addr": \)\(.*\)/\1"IP"/' configDelete.nft
sed -i -r 'N;s/.*(\n(\s*)"ip-addr":)/\2"name": "NAME",\1/;P;D' configDelete.nft
