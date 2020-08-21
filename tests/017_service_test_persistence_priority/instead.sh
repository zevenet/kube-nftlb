#!/bin/bash
kubectl apply -f . &>/dev/null

for file in priority*; do
	if [ -f $file ]; then
		sleep 10
		kubectl apply -f "$file" &>/dev/null
		curl --silent -H "Key: 12345" http://localhost:5555/farms/my-service--http >> configCreation.nft
		printf "\n\n" >> configCreation.nft
	fi
done

kubectl delete -f . &>/dev/null
sleep 10
curl --silent -H "Key: 12345" http://localhost:5555/farms/my-service--http > configDelete.nft

sed -i 's/\("virtual-addr": \)\(.*\)/\1"IP"/' configCreation.nft
sed -i 's/\("ip-addr": \)\(.*\)/\1"IP"/' configCreation.nft
sed -i -r 'N;s/.*(\n(\s*)"ip-addr":)/\2"name": "NAME",\1/;P;D' configCreation.nft