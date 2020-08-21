#!/bin/bash

# Test nodePort
# This script allows us to create a farm with type nodeport 
# Then it connects to an external device to make a request to the service
# for test information, see doc.md

NFTLB_KEY=$(grep 'NFTLB_KEY' ../../.env | sed 's/NFTLB_KEY=//')

# To perform the test you need a series of variables that you will find below. 
# Among which is the ssh key to make the connection, user and IP of the external client and your network interface to detect the IP of your machine
# In short, using the ssh command we connect to an external client and make a curl request to the service through the nodeport port.

keyssh="ssh_Key"
userClient="user_Client"
ipClient="ip_Client"
networkInterface="your_network_interface" # enviroment k8s

kubectl apply -f . &>/dev/null
sleep 10
curl --silent -H "Key: $NFTLB_KEY" http://localhost:5555/farms > configCreation.nft
# Check that the above parameters exist
if [[ -f "$keyssh" ]] && [ -n "$ipClient" ] && [ -n "$userClient" ] && [ -n "$userClient" ]
	then
		ipService=$(ifconfig $networkInterface | sed -En -e 's/.*inet ([0-9.]+).*/\1/p')
		ssh -i "$keyssh" "$userClient"@"$ipClient" curl --silent http://$ipService:32490 > requestNodeport.nft
	else
		echo "An error has occurred with the request, make sure you have configured the fields above correctly! abort test."
fi
kubectl delete -f . &>/dev/null
sleep 10
curl --silent -H "Key: $NFTLB_KEY" http://localhost:5555/farms > configDelete.nft

# Apply format
sed -i 's/\("virtual-addr": \)\(.*\)/\1"IP"/' configCreation.nft
sed -i 's/\("ip-addr": \)\(.*\)/\1"IP"/' configCreation.nft
sed -i -r 'N;s/.*(\n(\s*)"ip-addr":)/\2"name": "NAME",\1/;P;D' configCreation.nft
sed -i 's/\("virtual-addr": \)\(.*\)/\1"IP"/' configDelete.nft
sed -i 's/\("ip-addr": \)\(.*\)/\1"IP"/' configDelete.nft
sed -i -r 'N;s/.*(\n(\s*)"ip-addr":)/\2"name": "NAME",\1/;P;D' configDelete.nft
