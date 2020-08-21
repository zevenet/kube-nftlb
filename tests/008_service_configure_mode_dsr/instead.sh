#!/bin/bash

# Test DSR
# When configuring the DSR mode we have to make sure that the backends associated with the service have the service ip configured in their loopback network interface

# Create service & configure interfaces
kubectl apply -f service1.yaml &>/dev/null
kubectl apply -f deployment.yaml &>/dev/null
sleep 10
printf "DSR service creation \n" > configCreation.nft
curl --silent -H "Key: 12345" http://localhost:5555/farms/my-service--http >> configCreation.nft

results=( $(sed -n 's/.* "name": \([^ ]*\).*/\1/p' configCreation.nft) )
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
						printf "Interfaceloopback assign ip service \n" > configInterfaces.nft
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

# Change mode of service, check interfaces
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
curl --silent -H "Key: 12345" http://localhost:5555/farms/my-service--http >> configCreation.nft
 # Delete DSR Service
kubectl delete -f service2.yaml &>/dev/null
kubectl delete -f deployment.yaml &>/dev/null
sleep 10
curl --silent -H "Key: 12345" http://localhost:5555/farms/my-service--http > configDelete.nft
