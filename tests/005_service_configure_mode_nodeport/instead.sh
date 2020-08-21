#!/bin/bash

# Test nodePort
# This script allows us to create a farm with type nodeport 
# Then it connects to an external device to make a request to the service
# for test information, see doc.md

# To perform the test you need a series of variables that you will find below. 
# Among which is the ssh key to make the connection, user and IP of the external client and your network interface to detect the IP of your machine
# Through ssh, an http type request will be made

keyssh="keyNodeport"
userClient="david"
ipClient="192.168.10.182"
networkInterface="enp0s3"

kubectl apply -f . &>/dev/null
sleep 10
curl --silent -H "Key: 12345" http://localhost:5555/farms > configCreation.nft
if [[ -f "$keyssh" ]] && [ -n "$ipClient" ] && [ -n "$userClient" ] && [ -n "$userClient" ]
	then
		ipService=$(ifconfig $networkInterface | sed -En -e 's/.*inet ([0-9.]+).*/\1/p')
		ssh -i "$keyssh" "$userClient"@"$ipClient" curl --silent http://$ipService:32490 > requestNodeport.nft
	else
		echo "An error has occurred with the request, make sure you have configured the fields above correctly! abort test."
fi
sed -i 's/\("virtual-addr": \)\(.*\)/\1"IP"/' configCreation.nft
sed -i 's/\("ip-addr": \)\(.*\)/\1"IP"/' configCreation.nft
sed -i -r 'N;s/.*(\n(\s*)"ip-addr":)/\2"name": "NAME",\1/;P;D' configCreation.nft
kubectl delete -f . &>/dev/null
sleep 10
curl --silent -H "Key: 12345" http://localhost:5555/farms > configDelete.nft
