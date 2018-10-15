package pkg

import (
	kubernetes "k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
)

// GetClienset implements authentication to kube-nftlb.
// Stops the container if the authentication fails.
func GetClienset() *kubernetes.Clientset {
	// Get in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// Get Clientset for auth
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}
