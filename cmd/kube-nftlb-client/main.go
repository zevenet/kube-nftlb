package main

import (
	"fmt"
	"time"

	"github.com/zevenet/kube-nftlb/pkg/auth"
	"github.com/zevenet/kube-nftlb/pkg/config"
	"github.com/zevenet/kube-nftlb/pkg/logs"
	"github.com/zevenet/kube-nftlb/pkg/watchers"
	"k8s.io/apimachinery/pkg/util/wait"
)

func main() {
	// Make log channel before writing messages
	logChannel := make(chan string)
	levelLog := 0
	// Authentication: get access to the API
	clientset := auth.GetClienset(config.ClientCfgPath)
	go logs.PrintLogChannel(levelLog, fmt.Sprintf("%s", "Authentication successful"), logChannel)
	// Make lists of resources to be observed
	listWatchSvc := watchers.GetServiceListWatch(clientset)
	listWatchEndpoint := watchers.GetEndpointListWatch(clientset)
	go logs.PrintLogChannel(levelLog, fmt.Sprintf("%s", "Watchers ready"), logChannel)
	// Notify every change into logChannel based on every list watch
	serviceController := watchers.GetServiceController(listWatchSvc, logChannel, clientset)
	endpointController := watchers.GetEndpointController(listWatchEndpoint, logChannel, clientset)
	go logs.PrintLogChannel(levelLog, fmt.Sprintf("%s", "Controllers ready"), logChannel)
	// Event loop start, run them as background processes
	go serviceController.Run(wait.NeverStop)
	go logs.PrintLogChannel(levelLog, fmt.Sprintf("%s", "Service controller started"), logChannel)
	// We establish a waiting time for the creation of farms. This value is important or our farms will not be created correctly. Can be parameterized
	time.Sleep(config.ClientStartDelayTime)
	go endpointController.Run(wait.NeverStop)
	go logs.PrintLogChannel(levelLog, fmt.Sprintf("%s", "Endpoints controller started"), logChannel)
	// Print every message received from the channel
	go logs.PrintLogChannel(levelLog, fmt.Sprintf("%s", "Init finished"), logChannel)
	for message := range logChannel {
		fmt.Println(message)
	}
	// This line is unreachable: working as intended
}
