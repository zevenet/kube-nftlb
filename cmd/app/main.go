package main

import (
	"fmt"
	"time"

	auth "github.com/zevenet/kube-nftlb/pkg/auth"
	defaults "github.com/zevenet/kube-nftlb/pkg/defaults"
	watchers "github.com/zevenet/kube-nftlb/pkg/watchers"
	wait "k8s.io/apimachinery/pkg/util/wait"
)

func main() {
	// Read config values from the client (can be parameterized)
	cfg := defaults.Init()
	// Authentication: get access to the API
	clientset := auth.GetClienset(cfg.Global.KubeCfgPath)
	fmt.Println("Authentication successful...")
	// Make lists of resources to be observed
	listWatchSvc := watchers.GetServiceListWatch(clientset)
	listWatchEndpoint := watchers.GetEndpointListWatch(clientset)
	fmt.Println("Watchers ready...")
	// Make log channel before writing messages
	logChannel := make(chan string)
	// Notify every change into logChannel based on every list watch
	serviceController := watchers.GetServiceController(listWatchSvc, logChannel)
	endpointController := watchers.GetEndpointController(listWatchEndpoint, logChannel)
	fmt.Println("Controllers ready...")
	// Event loop start, run them as background processes
	go serviceController.Run(wait.NeverStop)
	fmt.Println("Service controller started")
	// We establish a waiting time for the creation of farms. This value is important or our farms will not be created correctly. Can be parameterized
	time.Sleep(time.Duration(cfg.Global.TimeStartApp) * time.Second)
	go endpointController.Run(wait.NeverStop)
	fmt.Println("Endpoints controller started")
	// Print every message received from the channel
	fmt.Println("Init finished")
	for message := range logChannel {
		fmt.Println(message)
	}
	// This line is unreachable: working as intended
}
