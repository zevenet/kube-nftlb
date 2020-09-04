package parser

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/zevenet/kube-nftlb/pkg/config"
	"github.com/zevenet/kube-nftlb/pkg/dsr"
	"github.com/zevenet/kube-nftlb/pkg/types"

	corev1 "k8s.io/api/core/v1"
)

// CreateServicePort
func CreateServicePort(service *corev1.Service, servicePort corev1.ServicePort, persistence string, persistTTL string, farms *[]types.Farm, wg *sync.WaitGroup) {
	// If we are creating a single service and it does not have a port name, we assign it a default name
	if servicePort.Name == "" {
		servicePort.Name = "default"
	}
	farmName := assignFarmNameService(service.Name, servicePort.Name)

	// Read the annotations collected in the "annotations" field of the service
	mode, scheduler, schedParam, helper, log, logPrefix := getAnnotations(service, farmName)

	// Get iface
	iface := ""
	if mode == "dsr" {
		iface = config.DockerInterfaceBridge
	}

	// Get farm struct with default values
	farm := types.Farm{
		Name:         farmName,
		Mode:         mode,
		Iface:        iface,
		IntraConnect: "on",
		Family:       findFamily(service),
		Helper:       helper,
		Log:          log,
		LogPrefix:    logPrefix,
		Persistence:  persistence,
		PersistTTL:   persistTTL,
		Protocol:     strings.ToLower(string(servicePort.Protocol)),
		Scheduler:    scheduler,
		SchedParam:   schedParam,
		State:        "up",
		VirtualAddr:  service.Spec.ClusterIP,
		VirtualPorts: fmt.Sprint(servicePort.Port),
		Backends:     make([]types.Backend, 0),
	}

	farmsMap[service.Name] = append(farmsMap[service.Name], farm.Name)
	*farms = append(*farms, farm)

	// Check if the service has DSR mode and append it to the DSR list
	if farm.Mode == "dsr" {
		go dsr.CreateInterfaceDsr(farm, service)
	} else if dsr.ExistsServiceDSR(farm.Name) && service.Spec.Type != "NodePort" {
		go dsr.DeleteInterfaceDsr(farm.Name)
	}

	// If the service type is NodePort, modify the name of the original service, its VP and its VIP
	if service.Spec.Type == "NodePort" || service.Spec.Type == "LoadBalancer" && servicePort.NodePort >= 0 {
		farm.Name = assignFarmNameNodePort(service.Name+"--"+servicePort.Name, "nodePort")
		farm.VirtualAddr = ""
		farm.VirtualPorts = fmt.Sprint(servicePort.NodePort)

		farmsMap[service.Name] = append(farmsMap[service.Name], farm.Name)
		*farms = append(*farms, farm)

		// Check if the service has DSR mode and append it to the DSR list
		if farm.Mode == "dsr" {
			go dsr.CreateInterfaceDsr(farm, service)
		} else if dsr.ExistsServiceDSR(farm.Name) && service.Spec.Type != "NodePort" {
			go dsr.DeleteInterfaceDsr(farm.Name)
		}
	}

	// Check if the service has externalIPs field
	for position, externalIPs := range service.Spec.ExternalIPs {
		farm.Name = assignFarmNameExternalIPs(service.Name+"--"+servicePort.Name, "externalIPs"+strconv.Itoa(position+1))
		farm.VirtualAddr = externalIPs

		farmsMap[service.Name] = append(farmsMap[service.Name], farm.Name)
		*farms = append(*farms, farm)
	}

	// Release this lock
	wg.Done()
}

// ServiceAsFarms
func ServiceAsFarms(service *corev1.Service) *types.Farms {
	// Make empty Farms struct where farms will be stored
	farms := make([]types.Farm, 0)

	// Gets persistence and Stickiness timeout in seconds
	persistence, persistTTL := getPersistence(service)
	findMaxConns(service)

	// Make wait group and process every service port
	var wg sync.WaitGroup
	wg.Add(len(service.Spec.Ports))
	for _, servicePort := range service.Spec.Ports {
		go CreateServicePort(service, servicePort, persistence, persistTTL, &farms, &wg)
	}

	// Wait until all locks are released
	wg.Wait()

	// Return the filled struct
	return &types.Farms{
		Farms: farms,
	}
}

// DeleteServiceFarms
func DeleteServiceFarms(service *corev1.Service, farmPathsChan chan<- string) {
	for _, farmName := range farmsMap[service.Name] {
		farmPathsChan <- fmt.Sprintf("farms/%s", farmName)
		go dsr.DeleteService(farmName)
	}
	close(farmPathsChan)
	delete(farmsMap, service.Name)
}

// DeleteMaxConnsService
func DeleteMaxConnsService(service *corev1.Service) {
	delete(maxConnsMap, service.Name)
}
