package parser

import (
	"fmt"
	"sync"

	"github.com/zevenet/kube-nftlb/pkg/dsr"
	"github.com/zevenet/kube-nftlb/pkg/types"

	corev1 "k8s.io/api/core/v1"
)

// ServiceAsFarms reads a Service and returns a filled Farms struct.
func ServiceAsFarms(service *corev1.Service) *types.Farms {
	// Make empty farms slice where farms will be stored
	farms := make([]types.Farm, 0)

	// Get persistence and persistTTL timeout in seconds
	persistence, persistTTL := getPersistence(service)

	// Find maxconns
	findMaxConns(service)

	// Make wait group and process every ServicePort
	var wg sync.WaitGroup
	wg.Add(len(service.Spec.Ports))
	for index := range service.Spec.Ports {
		// Can't reference the Service Port directly, because it only takes the first one and the rest is ignored
		go CreateServicePort(service, &service.Spec.Ports[index], persistence, persistTTL, &farms, &wg)
	}

	// Wait until all locks are released
	wg.Wait()

	// Return a Farms struct with the filled farms slice
	return &types.Farms{
		Farms: farms,
	}
}

// CreateServicePort appends farms parsed from a ServicePort to a Farm slice.
func CreateServicePort(service *corev1.Service, servicePort *corev1.ServicePort, persistence string, persistTTL string, farms *[]types.Farm, wg *sync.WaitGroup) {
	// Get a formatted farm name from Service and ServicePort
	farmName := FormatFarmName(service.Name, servicePort.Name)

	// Read the annotations collected in the "annotations" field of the service
	mode, scheduler, schedParam, helper, log, logPrefix := getAnnotations(service, farmName)

	// Get farm struct with default values
	farm := types.Farm{
		Name:         farmName,
		Mode:         mode,
		Helper:       helper,
		Log:          log,
		LogPrefix:    logPrefix,
		Persistence:  persistence,
		PersistTTL:   persistTTL,
		Scheduler:    scheduler,
		SchedParam:   schedParam,
		VirtualAddr:  service.Spec.ClusterIP,
		VirtualPorts: findVirtualPorts(servicePort),
		Protocol:     findProtocol(servicePort),
		Family:       findFamily(service),
		Iface:        findIface(mode),
		IntraConnect: "on",
		State:        "up",
		Backends:     make([]types.Backend, 0),
	}
	*farms = append(*farms, farm)

	go func() {
		// Append farm name to the slice of farm names based on the Service name
		farmsPerServiceMap[service.Name] = append(farmsPerServiceMap[service.Name], farm.Name)

		// Check DSR mode
		if farm.Mode == "dsr" {
			dsr.CreateInterfaceDsr(farm, service)
		}
	}()

	// If the Service type is NodePort, modify the farm name, its VP and its VIP
	if service.Spec.Type == "NodePort" || service.Spec.Type == "LoadBalancer" && servicePort.NodePort >= 0 {
		farm.Name = FormatNodePortFarmName(service.Name, servicePort.Name)
		farm.VirtualAddr = ""
		farm.VirtualPorts = findVirtualPortsNodePort(servicePort)
		*farms = append(*farms, farm)

		go func() {
			// Append NodePort farm name to the slice of farm names and NodePort based on the Service name
			farmsPerServiceMap[service.Name] = append(farmsPerServiceMap[service.Name], farm.Name)

			// Check DSR mode
			if farm.Mode == "dsr" {
				dsr.CreateInterfaceDsr(farm, service)
			}
		}()
	}

	// Make a farm for every externalIP found in the Service
	for index, externalIP := range service.Spec.ExternalIPs {
		farm.Name = FormatExternalIPFarmName(service.Name, servicePort.Name, index+1)
		farm.VirtualAddr = externalIP
		*farms = append(*farms, farm)

		go func() {
			// Append ExternalIP farm name to the slice of farm names based on the Service name
			farmsPerExternalIPResourceMap[service.Name] = append(farmsPerExternalIPResourceMap[service.Name], farm.Name)
		}()
	}

	// Release lock
	wg.Done()
}

// DeleteServiceFarms sends farm paths through a channel to the controller. The controller then deletes every farm.
func DeleteServiceFarms(service *corev1.Service, farmPathsChan chan<- string) {
	for _, farmName := range farmsPerServiceMap[service.Name] {
		farmPathsChan <- fmt.Sprintf("farms/%s", farmName)
		go dsr.DeleteService(farmName)
	}
	farmsPerServiceMap[service.Name] = []string{}

	for _, farmName := range farmsPerExternalIPResourceMap[service.Name] {
		farmPathsChan <- fmt.Sprintf("farms/%s", farmName)
	}
	farmsPerExternalIPResourceMap[service.Name] = []string{}

	close(farmPathsChan)
}

// DeleteMaxConnsService deletes the key:value pair from maxConnsMap.
func DeleteMaxConnsService(service *corev1.Service) {
	delete(maxConnsMap, service.Name)
}
