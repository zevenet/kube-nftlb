package auth

import (
	"flag"
	"fmt"

	"github.com/zevenet/kube-nftlb/pkg/config"
	"github.com/zevenet/kube-nftlb/pkg/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var clienset = authenticate(config.ClientCfgPath)

// GetClientset
func GetClientset() *kubernetes.Clientset {
	return clienset
}

// authenticate implements authentication to kube-nftlb. Stops the container if the authentication fails.
func authenticate(cfg string) *kubernetes.Clientset {
	// Parse command line flags
	kubeconfig := flag.String("kubeconfig", cfg, cfg)
	flag.Parse()

	// Collects the path of the configuration file and generates the necessary configuration to then create the client
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	// Create the clientset, based on the previous configuration
	clienset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	go log.WriteLog(0, fmt.Sprintf("%s", "Authentication successful"))

	return clienset
}
