package main

import (
	"fmt"

	watchers "github.com/zevenet/kube-nftlb/pkg/watchers"
	auth "github.com/zevenet/kube-nftlb/pkg/auth"
	wait "k8s.io/apimachinery/pkg/util/wait"
)

func main() {
	// Authentication: get access to the API
	clientset := auth.GetClienset()
	// Make lists of resources to be observed
	listWatchSvc := watchers.GetServiceListWatch(clientset)
	listWatchEndpoint := watchers.GetEndpointListWatch(clientset)
	// Make log channel before writing messages
	logChannel := make(chan string)
	// Notify every change into logChannel based on every list watch
	serviceController := watchers.GetServiceController(listWatchSvc, logChannel)
	endpointController := watchers.GetEndpointController(listWatchEndpoint, logChannel)
	// Event loop start, run them as background processes
	go serviceController.Run(wait.NeverStop)
	go endpointController.Run(wait.NeverStop)
	// Print every message received from the channel
	for message := range logChannel {
		fmt.Println(message)
	}
	// This line is unreachable: working as intended
}
