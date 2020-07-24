package auth

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// GetClienset implements authentication to kube-nftlb. Stops the container if the authentication fails.
func GetClienset(cfg string) *kubernetes.Clientset {
	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", cfg, cfg)
	flag.Parse()
	// collects the path of the configuration file and generates the necessary configuration to then create the client
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// create the clientset, based on the previous configuration
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}
