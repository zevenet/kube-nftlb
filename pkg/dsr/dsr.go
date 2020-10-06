package dsr

// TODO Adapt DSR mode to new Addresses nftlb object

/*
import (
	"context"
	"fmt"
	"strings"

	"github.com/zevenet/kube-nftlb/pkg/auth"
	"github.com/zevenet/kube-nftlb/pkg/types"

	dockerTypes "github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	mapServiceDsr = make(map[string]*types.DSR)
	clientset     = auth.GetClientset()
)

// CreateInterfaceDsr looks for the label fields within our YAML configuration file and lets us to identify which deployment is assigned to our service.
func CreateInterfaceDsr(farm types.Farm, service *corev1.Service) {
	if !ExistsServiceDSR(farm.Name) {
		mapServiceDsr[farm.Name] = new(types.DSR)
	}
	mapServiceDsr[farm.Name].VirtualAddr = farm.VirtualAddr
	mapServiceDsr[farm.Name].VirtualPorts = farm.VirtualPorts

	// Get service label and check if it exists
	labelService, ok := service.ObjectMeta.Labels["app"]
	if !ok {
		return
	}

	// Get pod list
	podList, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		// Don't continue if there are errors
		return
	} else if len(podList.Items) == 0 {
		// Don't continue if there aren't pods
		return
	}

	// If the first pods has "app=" label, match it with the deployment label
	labelDeployment, ok := podList.Items[0].ObjectMeta.Labels["app"]
	if !ok {
		// Don't continue if there are errors
		return
	} else if labelService != labelDeployment {
		// Don't continue if both labels don't match
		return
	}

	// If both labels match, add every pod (backend) name to the interface
	for _, pod := range podList.Items {
		// Configure the loopback interface with the IP of the service
		go AddInterfaceDsr(farm.Name, pod.ObjectMeta.Name, service.Spec.ClusterIP)
	}
}

// AddInterfaceDsr
func AddInterfaceDsr(farmName string, backendName string, virtualAddr string) {
	// What we do in this function is use the "dockerclient" to get the dockerUID of our containers.
	// This is necessary to locate the specific container on which we are going to configure the interface and
	// then on which we are going to execute the commands. For example: Search among all the pods, those that
	// are in the "default" namespace for those that are called as one of our deployments (this is launched once for each of them)
	pod, err := clientset.CoreV1().Pods("default").Get(context.TODO(), backendName, metav1.GetOptions{})
	if err != nil {
		// Don't continue if there are errors
		return
	} else if len(pod.Status.ContainerStatuses) == 0 {
		// Don't continue if pod doesn't have any status
		return
	}

	// Get UID list
	dockerUID := pod.Status.ContainerStatuses[0].ContainerID
	UIDs := strings.SplitAfter(dockerUID, "docker://")
	if len(UIDs) == 0 {
		// Don't continue if UID list is empty
		return
	}

	// We specify the configuration to apply on the container. The CMD field are the commands to apply.
	// The users field are the privileges with which these commands will be executed.
	// Network configuration needs to be run in root mode.
	execConfig := dockerTypes.ExecConfig{
		AttachStderr: true,
		AttachStdin:  true,
		AttachStdout: true,
		Cmd:          []string{"/bin/sh", "-c", "ip ad add " + virtualAddr + "/32 dev lo"},
		Tty:          true,
		Detach:       false,
		Privileged:   true,
		User:         "root",
		WorkingDir:   "/",
	}

	execAttachConfig := dockerTypes.ExecStartCheck{
		Detach: false,
		Tty:    true,
	}

	// Use the docker client to make requests to the docker rest api, from it we get the UID
	// of the pod that we want to apply DSR
	cli, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv)
	if err != nil {
		panic(err)
	}

	// Link the configuration to the container where we want to apply said configuration
	// (once linked, it is only necessary to launch the "Start" command)
	exec, err := cli.ContainerExecCreate(context.TODO(), UIDs[1], execConfig)
	if err != nil {
		panic(err)
	}

	err = cli.ContainerExecStart(context.TODO(), exec.ID, execAttachConfig)
	if err != nil {
		panic(err)
	}

	// Store the UID of this container to use it later
	mapServiceDsr[farmName].DockerUID = append(mapServiceDsr[farmName].DockerUID, UIDs[1])
}

// DeleteInterfaceDsr
func DeleteInterfaceDsr(farmName string) {
	// Delete the configuration of the loopback interface of our deployments
	for UID := range mapServiceDsr[farmName].DockerUID {
		// Get configs
		execConfig := dockerTypes.ExecConfig{
			AttachStderr: true,
			AttachStdin:  true,
			AttachStdout: true,
			Cmd:          []string{"/bin/sh", "-c", "ip ad del " + mapServiceDsr[farmName].VirtualAddr + "/32 dev lo"},
			Tty:          true,
			Detach:       false,
			Privileged:   true,
			User:         "root",
			WorkingDir:   "/",
		}

		execAttachConfig := dockerTypes.ExecStartCheck{
			Detach: false,
			Tty:    true,
		}

		// New Docker client
		cli, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv)
		if err != nil {
			panic(err)
		}

		// Create container with exec config
		exec, err := cli.ContainerExecCreate(context.TODO(), fmt.Sprintf("%s", mapServiceDsr[farmName].DockerUID[UID]), execConfig)
		if err != nil {
			panic(err)
		}

		// Start container with exec attach config
		err = cli.ContainerExecStart(context.TODO(), exec.ID, execAttachConfig)
		if err != nil {
			panic(err)
		}
	}

	// Remove this service from the DSR map
	go DeleteService(farmName)
}

// DeleteService
func DeleteService(farmName string) {
	if ExistsServiceDSR(farmName) {
		delete(mapServiceDsr, farmName)
	}
}

// ExistsServiceDSR
func ExistsServiceDSR(farmName string) bool {
	_, exists := mapServiceDsr[farmName]
	return exists
}

func GetVirtualAddr(farmName string) string {
	return mapServiceDsr[farmName].VirtualAddr
}

func GetVirtualPorts(farmName string) string {
	return mapServiceDsr[farmName].VirtualPorts
}
*/
