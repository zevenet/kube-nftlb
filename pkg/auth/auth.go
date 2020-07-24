package auth

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

// GetClienset implements authentication to kube-nftlb. Stops the container if the authentication fails.
func GetClienset() *kubernetes.Clientset {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join("/", "var", "config-kubernetes/admin.conf"), "/var/config-kubernetes/admin.conf")
	} else {

		kubeconfig = flag.String("kubeconfig", "", "/var/config-kubernetes/admin.conf")
	}
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

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows users
}
