package dsr

import (
	"context"
	"fmt"
	"strings"

	dockerClient "github.com/docker/docker/client"
	dockerTypes "github.com/docker/docker/api/types"
	kubernetes "k8s.io/client-go/kubernetes"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "github.com/zevenet/kube-nftlb/pkg/types"
)

// Every time the main map is updated, this variable does the same. We then use it to retrieve the data contained within the map.
var mapServiceDsr = make(map[string]*types.DsrStruct)

func GetDsrArray() map[string]*types.DsrStruct {
	return mapServiceDsr
}

func CreateInterfaceDsr(farmName string, service *v1.Service, clientset *kubernetes.Clientset, serviceDsr map[string]*types.DsrStruct) map[string]*types.DsrStruct  {
	// What we do in this function is look for the label fields within our yaml configuration file. This allows us to identify which deployment is assigned to our service.
	// Once located we can configure the loopback interface with the ip of the service.
	virtualAddr := service.Spec.ClusterIP
	if _, ok := service.ObjectMeta.Labels["app"]; ok {
		labelService := service.ObjectMeta.Labels["app"]
		objPod, _ := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
		if len(objPod.Items) >= 1 {
			if _, ok := objPod.Items[0].ObjectMeta.Labels["app"]; ok {
				labelDeployment := objPod.Items[0].ObjectMeta.Labels["app"]
				if labelService == labelDeployment {
					for _, objectMeta := range objPod.Items {
						backendName := objectMeta.ObjectMeta.Name
						// Call the function that configures the interface of that specific deployment
						serviceDsr = AddInterfaceDsr(clientset, farmName, backendName, virtualAddr, serviceDsr)
					}
				}
			}
		}
	}
	mapServiceDsr = serviceDsr
	return serviceDsr
}

func AddInterfaceDsr(clientset *kubernetes.Clientset, farmName string, backendName string, virtualAddr string, serviceDsr map[string]*types.DsrStruct) (map[string]*types.DsrStruct) {
	// What we do in this function is use the "dockerclient" to get the dockerUid of our containers.
	// This is necessary to locate the specific container on which we are going to configure the interface and then on which we are going to execute the commands.
	// For example: Search among all the pods, those that are in the "default" namespace for those that are called as one of our deployments (this is launched once for each of them)
	objContainer, _ := clientset.CoreV1().Pods("default").Get(context.TODO(), backendName, metav1.GetOptions{})
	if len(objContainer.Status.ContainerStatuses) >= 1 {
		dockerUid := objContainer.Status.ContainerStatuses[0].ContainerID
		uid := strings.SplitAfter(dockerUid, "docker://")
		if len(uid) >= 1 {
			// We use the docker client to make requests to the docker rest api. From it we get the UID of the pod that we want to apply DSR.
			cli, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv)
			if err != nil {
				panic(err)
			}
			// We specify the configuration to apply on the container. The CMD field are the commands to apply. The users field are the privileges with which these commands will be executed.
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
			// We link the configuration to the container where we want to apply said configuration. Once linked, it is only necessary to launch the "Start" command
			exec, err := cli.ContainerExecCreate(context.TODO(), uid[1], execConfig)
			if err != nil {
				panic(err)
			}
			execAttachConfig := dockerTypes.ExecStartCheck{
				Detach: false,
				Tty:    true,
			}
			err = cli.ContainerExecStart(context.TODO(), exec.ID, execAttachConfig)
			if err != nil {
				panic(err)
			}
			// We save the uid of our container to use it later
			serviceDsr[farmName].DockerUid = append(serviceDsr[farmName].DockerUid, uid[1])
		}
	}
	mapServiceDsr = serviceDsr
	return serviceDsr
}

func DeleteInterfaceDsr(farmName string, serviceDsr map[string]*types.DsrStruct) map[string]*types.DsrStruct {
	// In this function we take care of deleting the configuration of the loopback interface of our deployments and we also remove the service from the DSR list
	for uid := range serviceDsr[farmName].DockerUid {
		cli, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv)
		if err != nil {
			panic(err)
		}
		execConfig := dockerTypes.ExecConfig{
			AttachStderr: true,
			AttachStdin:  true,
			AttachStdout: true,
			Cmd:          []string{"/bin/sh", "-c", "ip ad del " + serviceDsr[farmName].VirtualAddr + "/32 dev lo"},
			Tty:          true,
			Detach:       false,
			Privileged:   true,
			User:         "root",
			WorkingDir:   "/",
		}

		exec, err := cli.ContainerExecCreate(context.TODO(), fmt.Sprintf("%s", serviceDsr[farmName].DockerUid[uid]), execConfig)
		if err != nil {
			panic(err)
		}
		execAttachConfig := dockerTypes.ExecStartCheck{
			Detach: false,
			Tty:    true,
		}
		err = cli.ContainerExecStart(context.TODO(), exec.ID, execAttachConfig)
		if err != nil {
			panic(err)
		}
	}
	serviceDsr = DeleteServiceDsr(farmName, serviceDsr)
	mapServiceDsr = serviceDsr
	return serviceDsr
}

func DeleteServiceDsr(key string, serviceDsr map[string]*types.DsrStruct) map[string]*types.DsrStruct {
	delete(serviceDsr, key)
	mapServiceDsr = serviceDsr
	return serviceDsr
}