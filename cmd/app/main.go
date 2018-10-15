package main

import (
	"fmt"

	auth "github.com/zevenet/kube-nftlb/pkg/auth"
	svc "github.com/zevenet/kube-nftlb/pkg/svc"
	wait "k8s.io/apimachinery/pkg/util/wait"
)

func main() {
	// Authentication: get access to the API
	clientset := auth.GetClienset()
	// Make list of resources (every Service) to be observed
	listWatch := svc.GetServiceListWatch(clientset)
	// Make log channel before writing messages
	logChannel := make(chan string)
	// Notify every change into logChannel based on listWatch
	serviceController := svc.GetServiceController(listWatch, logChannel)
	// Event loop start, run it as background process
	go serviceController.Run(wait.NeverStop)
	// Print every message received from the channel
	for message := range logChannel {
		fmt.Println(message)
	}
	// This line is unreachable: working as intended
}
