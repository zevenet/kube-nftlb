#!/bin/bash

# Test DSR
# When configuring the DSR mode we have to make sure that the backends associated with the service have the service ip configured in their loopback network interface

NFTLB_KEY=$(grep 'NFTLB_KEY' ../../.env | sed 's/NFTLB_KEY=//')

# Create service & configure interfaces
kubectl apply -f service1.yaml &>/dev/null
kubectl apply -f deployment.yaml &>/dev/null
sleep 10
printf "DSR service creation \n" > configCreation.nft
curl --silent -H "Key: $NFTLB_KEY" http://localhost:5555/farms/my-service--http >> configCreation.nft
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
						printf "Interfaceloopback assign ip service \n" > configInterfaces.nft
						# We show the loopback network interface of our deployments
						kubectl exec -it "$temp" -n "default" -- ip a show lo >> configInterfaces.nft
						printf "\n" >> configInterfaces.nft
				elif ((n == 2))
					then
						kubectl exec -it "$temp" -n "default" -- ip a show lo >> configInterfaces.nft 
						printf "\n" >> configInterfaces.nft
				fi
		fi
	((n++))
done

# Once we have verified that the interface has been configured, we use the second configuration file yaml to also verify that by eliminating the DSR mode, 
# the interface configuration is also eliminated. We use the same procedure as before.
printf "\n DSR service change mode \n" >> configCreation.nft
kubectl apply -f service2.yaml &>/dev/null
sleep 10
n=0
for container in "${results[@]}"
	do
		if ((n != 0)) 
			then
				sleep 5
				temp=$(echo "${container::-2}")
				temp=$(echo "${temp:1}")
				if ((n == 1))
					then
						printf "Interfaceloopback deleted ip service \n" >> configInterfaces.nft
						kubectl exec -it "$temp" -n "default" -- ip a show lo >> configInterfaces.nft
						printf "\n" >> configInterfaces.nft
				elif ((n == 2))
					then
						kubectl exec -it "$temp" -n "default" -- ip a show lo >> configInterfaces.nft 
						printf "\n" >> configInterfaces.nft
				fi
		fi
	((n++))
done
curl --silent -H "Key: $NFTLB_KEY" http://localhost:5555/farms/my-service--http >> configCreation.nft
 # Delete DSR Service
kubectl delete -f service2.yaml &>/dev/null
kubectl delete -f deployment.yaml &>/dev/null
sleep 10
curl --silent -H "Key: $NFTLB_KEY" http://localhost:5555/farms/my-service--http > configDelete.nft
# Apply format
sed -i 's/\("virtual-addr": \)\(.*\)/\1"IP"/' configCreation.nft
sed -i 's/\("ip-addr": \)\(.*\)/\1"IP"/' configCreation.nft
sed -i -r 'N;s/.*(\n(\s*)"ip-addr":)/\2"name": "NAME",\1/;P;D' configCreation.nft
sed -i 's/\("virtual-addr": \)\(.*\)/\1"IP"/' configDelete.nft
sed -i 's/\("ip-addr": \)\(.*\)/\1"IP"/' configDelete.nft
sed -i -r 'N;s/.*(\n(\s*)"ip-addr":)/\2"name": "NAME",\1/;P;D' configDelete.nft
sed -i 's/.*scope global lo.*/       IP_SERVICE/' configInterfaces.nft
