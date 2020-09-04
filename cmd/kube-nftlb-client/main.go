package main

import (
	"time"

	"github.com/zevenet/kube-nftlb/pkg/auth"
	"github.com/zevenet/kube-nftlb/pkg/config"
	"github.com/zevenet/kube-nftlb/pkg/controller"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
)

func main() {
	// Authentication: get access to the API
	clientset := auth.GetClientset()

	// Get controllers
	controllers := []cache.Controller{
		controller.NewServiceController(clientset),
		controller.NewEndpointsController(clientset),
		controller.NewNetworkPolicyController(clientset),
	}

	// Run controllers as background processes
	for _, controller := range controllers {
		go controller.Run(wait.NeverStop)
		time.Sleep(config.ClientStartDelayTime)
	}

	select {}
	// This line is unreachable: working as intended
}
