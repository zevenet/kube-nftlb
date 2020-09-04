package json

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zevenet/kube-nftlb/pkg/config"
	"github.com/zevenet/kube-nftlb/pkg/dsr"
	"github.com/zevenet/kube-nftlb/pkg/types"

	configFarm "github.com/zevenet/kube-nftlb/pkg/farms"
	corev1 "k8s.io/api/core/v1"
)

var farmsMap = map[string][]string{}

// CreateServicePort
func CreateServicePort(service *corev1.Service, servicePort corev1.ServicePort, persistence string, persistTTL string) []types.Farm {
	servicePortFarms := make([]types.Farm, 0)

	// If we are creating a single service and it does not have a port name, we assign it a default name
	if servicePort.Name == "" {
		servicePort.Name = "default"
	}
	farmName := configFarm.AssignFarmNameService(service.ObjectMeta.Name, servicePort.Name)

	// Read the annotations collected in the "annotations" field of the service
	mode, scheduler, schedulerParam, helper, log, logPrefix := getAnnotations(service, farmName)

	// Get iface
	iface := ""
	if mode == "dsr" {
		iface = config.DockerInterfaceBridge
	}

	// Get farm struct with default values
	farm := types.Farm{
		Name:           farmName,
		Mode:           mode,
		Iface:          iface,
		Intraconnect:   "on",
		Family:         findFamily(service),
		Helper:         helper,
		Log:            log,
		LogPrefix:      logPrefix,
		Persistence:    persistence,
		PersistTTL:     persistTTL,
		Protocol:       strings.ToLower(string(servicePort.Protocol)),
		Scheduler:      scheduler,
		SchedulerParam: schedulerParam,
		State:          "up",
		VirtualAddr:    service.Spec.ClusterIP,
		VirtualPorts:   fmt.Sprint(servicePort.Port),
		Backends:       make([]types.Backend, 0),
	}

	farmsMap[service.ObjectMeta.Name] = append(farmsMap[service.ObjectMeta.Name], farm.Name)
	servicePortFarms = append(servicePortFarms, farm)

	// Check if the service has DSR mode and append it to the DSR list
	if farm.Mode == "dsr" {
		go dsr.CreateInterfaceDsr(farm, service)
	} else if dsr.ExistsServiceDSR(farm.Name) && service.Spec.Type != "NodePort" {
		go dsr.DeleteInterfaceDsr(farm.Name)
	}

	// If the service type is NodePort, modify the name of the original service, its VP and its VIP
	if service.Spec.Type == "NodePort" || service.Spec.Type == "LoadBalancer" && servicePort.NodePort >= 0 {
		farm.Name = configFarm.AssignFarmNameNodePort(service.ObjectMeta.Name+"--"+servicePort.Name, "nodePort")
		farm.VirtualAddr = ""
		farm.VirtualPorts = fmt.Sprint(servicePort.NodePort)

		farmsMap[service.ObjectMeta.Name] = append(farmsMap[service.ObjectMeta.Name], farm.Name)
		servicePortFarms = append(servicePortFarms, farm)

		// Check if the service has DSR mode and append it to the DSR list
		if farm.Mode == "dsr" {
			go dsr.CreateInterfaceDsr(farm, service)
		} else if dsr.ExistsServiceDSR(farm.Name) && service.Spec.Type != "NodePort" {
			go dsr.DeleteInterfaceDsr(farm.Name)
		}
	}

	// Check if the service has externalIPs field
	for position, externalIPs := range service.Spec.ExternalIPs {
		farm.Name = configFarm.AssignFarmNameExternalIPs(service.ObjectMeta.Name+"--"+servicePort.Name, "externalIPs"+strconv.Itoa(position+1))
		farm.VirtualAddr = externalIPs

		farmsMap[service.ObjectMeta.Name] = append(farmsMap[service.ObjectMeta.Name], farm.Name)
		servicePortFarms = append(servicePortFarms, farm)
	}

	fmt.Println("Finished serviceport " + servicePort.Name + " from " + farm.Name)

	return servicePortFarms
}

// ParseServiceAsFarms
func ParseServiceAsFarms(service *corev1.Service) types.Farms {
	// Make empty Farms struct where farms will be stored
	farms := types.Farms{
		Farms: make([]types.Farm, 0),
	}

	// Gets persistence and Stickiness timeout in seconds
	persistence, persistTTL := getPersistence(service)
	findMaxConns(service)

	// When creating services we can create several from the same yaml configuration file
	// For this we take into account the port field of the yaml configuration file. We create a service for each name field in ports
	for _, servicePort := range service.Spec.Ports {
		fmt.Println("Added serviceport " + servicePort.Name)
		farms.Farms = append(farms.Farms, CreateServicePort(service, servicePort, persistence, persistTTL)...)
	}

	// Returns the filled struct
	return farms
}

//
func DeleteServiceFarms(service *corev1.Service, farmPathsChan chan<- string) {
	for _, farmName := range farmsMap[service.ObjectMeta.Name] {
		farmPathsChan <- fmt.Sprintf("farms/%s", farmName)
		go dsr.DeleteService(farmName)
	}
	close(farmPathsChan)
	go delete(farmsMap, service.ObjectMeta.Name)
}

// DeleteMaxConnsMap
func DeleteMaxConnsMap() {
	maxConnsMap = map[string]string{}
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}

	// Not found
	return -1
}
