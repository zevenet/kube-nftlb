package dsr

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/zevenet/kube-nftlb/pkg/types"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// Map [farm (name)] to { DSR object }
	farmsDSR = make(map[string]types.DSR)
)

// IsEnabled ...
func IsEnabled(farm *types.Farm) bool {
	_, exists := farmsDSR[farm.Name]
	return exists
}

// Enable ...
func Enable(farm *types.Farm) {
	farmsDSR[farm.Name] = types.DSR{
		DockerUIDs:   make([]string, 0),
		AddressesIPs: make([]string, len(farm.Addresses)),
	}

	for index, address := range farm.Addresses {
		farmsDSR[farm.Name].AddressesIPs[index] = address.IPAddr
	}
}

// Disable ...
func Disable(farm *types.Farm) {
	deleteInterfaces(farm.Name)
	delete(farmsDSR, farm.Name)
}

// createInterface looks for the label fields within our YAML configuration file and lets us to identify which deployment is assigned to our service.
func createInterfaces(farm *types.Farm, service *corev1.Service) {
	podList := getPodList(service)
	if podList == nil {
		return
	}

	// If both labels match, add every pod (backend) name to the interface
	for _, pod := range podList {
		// Configure the loopback interface with the IP of the service
		addPodInterface(farm.Name, pod.ObjectMeta.Name)
	}
}

// addPodInterface
func addPodInterface(farmName string, backendName string) error {
	UID, err := getPodUID(backendName)
	if err != nil {
		return err
	}

	// Store the UID of this container to use it later
	farmDSR := farmsDSR[farmName]
	farmDSR.DockerUIDs = append(farmDSR.DockerUIDs, UID)
	farmsDSR[farmName] = farmDSR

	for _, virtualIP := range farmsDSR[farmName].AddressesIPs {
		if err := dockerCmdRun("add", virtualIP, UID); err != nil {
			return err
		}
	}

	return nil
}

// DeleteInterfaces
func deleteInterfaces(farmName string) {
	// Delete the configuration of the loopback interface of our deployments
	for indexUID := range farmsDSR[farmName].DockerUIDs {
		for _, virtualIP := range farmsDSR[farmName].AddressesIPs {
			if err := dockerCmdRun("del", virtualIP, farmsDSR[farmName].DockerUIDs[indexUID]); err != nil {
				panic(err)
			}
		}
	}
}

func getPodList(service *corev1.Service) []corev1.Pod {
	// Get service label and check if it exists
	labelService, ok := service.ObjectMeta.Labels["app"]
	if !ok {
		return nil
	}

	// Get pod list
	podList, err := clientset.CoreV1().Pods(corev1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		// Don't continue if there are errors
		return nil
	} else if len(podList.Items) == 0 {
		// Don't continue if there aren't pods
		return nil
	}

	// If the first pods has "app=" label, match it with the deployment label
	labelDeployment, ok := podList.Items[0].ObjectMeta.Labels["app"]
	if !ok {
		// Don't continue if there are errors
		return nil
	} else if labelService != labelDeployment {
		// Don't continue if both labels don't match
		return nil
	}

	return podList.Items
}

func getPodUID(backendName string) (string, error) {
	pod, err := clientset.CoreV1().Pods(corev1.NamespaceAll).Get(context.TODO(), backendName, metav1.GetOptions{})
	if err != nil {
		// Don't continue if there are errors
		return "", err
	} else if len(pod.Status.ContainerStatuses) == 0 {
		// Don't continue if pod doesn't have any status
		return "", errors.New("pod.Status.ContainerStatuses is empty")
	}

	splittedUID := strings.SplitAfter(pod.Status.ContainerStatuses[0].ContainerID, "docker://")
	if len(splittedUID) != 2 {
		// Don't continue if the splitted UID doesn't have 2 elements
		return "", fmt.Errorf("Invalid length after splitting %s from \"docker://\"", pod.Status.ContainerStatuses[0].ContainerID)
	}

	// Return the second element from ["docker://", UID]
	return splittedUID[1], nil
}
