package http

import (
	"os"
)

var (
	// Hostname is discovered dynamically.
	Hostname string

	// BadNames is a name list of pods/services that shouldn't be doing any requests (they have invalid data).
	BadNames []string
)

func init() {
	var err error

	// Discover hostname
	Hostname, err = os.Hostname()
	if err != nil {
		panic(err)
	}

	// Set blacklist
	BadNames = []string{"kube-controller-manager", "kube-scheduler", "kube-scheduler-" + Hostname, "kube-controller-manager-" + Hostname, "k8s.io-minikube-hostpath"}
}
