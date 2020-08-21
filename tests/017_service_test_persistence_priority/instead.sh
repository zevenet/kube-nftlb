#!/bin/bash

# Test Persistance
# This script allows us to create a farm with two associated backends. 
# Then add the mark field in the backends and replace the index.html of each of our backends so that when we make an http request we know which backend is responding.
# for test information, see doc.md

NFTLB_KEY=$(grep 'NFTLB_KEY' ../../.env | sed 's/NFTLB_KEY=//')

kubectl apply -f . &>/dev/null
sleep 10
curl --silent -H "Key: $NFTLB_KEY" http://localhost:5555/farms/my-service--http > configCreation.nft
# Gets all the backend names and loops through them one by one 
results=( $(sed -n 's/.* "name": \([^ ]*\).*/\1/p' configCreation.nft) )
n=0
for container in "${results[@]}"
do
	if ((n != 0)) 
	then
		sleep 5
		# Remove unnecessary characters from text strings
		temp=$(echo "${container::-2}")
		temp=$(echo "${temp:1}")
		if ((n == 1))
		then
			# Currently to test persistence you have to add marks to the backends (In the future this will be automatic)
			# Each marks must be different for each backend. Ej 0x2 & 0x4
			jsonPath='{"farms" : [ { "name" : "my-service--http", "backends" : [ { "name" : "%s", "mark" : "0x00000002" } ] } ] }\n'
			jsonCurl=$(printf "$jsonPath" "$temp")
			curl --silent -H "Key: $NFTLB_KEY" -X POST http://localhost:5555/farms -d "$jsonCurl" &>/dev/null
		elif ((n == 2))
		then
			jsonPath='{"farms" : [ { "name" : "my-service--http", "backends" : [ { "name" : "%s", "mark" : "0x00000004" } ] } ] }\n'
			jsonCurl=$(printf "$jsonPath" "$temp")
			curl --silent -H "Key: $NFTLB_KEY" -X POST http://localhost:5555/farms -d "$jsonCurl" &>/dev/null
		fi
		# Modify the index.html to know which backend is responding to us
		kubectl exec -it "$temp" -n "default" -- /bin/sh -c "echo backend$n > /usr/share/nginx/html/index.html"
	fi
	((n++))
done
curl --silent -H "Key: $NFTLB_KEY" http://localhost:5555/farms/my-service--http > configCreation.nft
sed -i 's/\("virtual-addr": \)\(.*\)/\1"IP"/' configCreation.nft
sed -i 's/\("ip-addr": \)\(.*\)/\1"IP"/' configCreation.nft
sed -i -r 'N;s/.*(\n(\s*)"ip-addr":)/\2"name": "NAME",\1/;P;D' configCreation.nft
ipService=$(kubectl get service/my-service -o jsonpath='{.spec.clusterIP}')
curl --silent http://$ipService:8080 > requestService.nft
curl --silent http://$ipService:8080 >> requestService.nft
curl --silent http://$ipService:8080 >> requestService.nft
curl --silent http://$ipService:8080 >> requestService.nft
kubectl delete -f . &>/dev/null
sleep 10
curl --silent -H "Key: $NFTLB_KEY" http://localhost:5555/farms/my-service--http > configDelete.nft
