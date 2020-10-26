package parser

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/zevenet/kube-nftlb/pkg/config"
	"github.com/zevenet/kube-nftlb/pkg/dsr"
	"github.com/zevenet/kube-nftlb/pkg/types"

	corev1 "k8s.io/api/core/v1"
)

// ServiceAsPaths sends farm and addresses paths through a channel to the controller. The controller then sends a
// DELETE request to nftlb for every path.
func ServiceAsPaths(service *corev1.Service, pathChan chan<- string) {
	for _, farmName := range farmsPerService[service.Name] {
		// Send farm path to the controller
		pathChan <- fmt.Sprintf("farms/%s", farmName)

		// Send addresses paths to the controller
		for _, addressName := range addressesPerFarm[farmName] {
			pathChan <- fmt.Sprintf("addresses/%s", addressName)
		}

		// Remove from memory addresses names mapped to this farm
		delete(addressesPerFarm, farmName)
	}

	// Remove from memory farm names mapped to this Service
	delete(farmsPerService, service.Name)

	close(pathChan)
}

// ServiceAsNftlb analyzes a Service and returns a filled Nftlb struct.
func ServiceAsNftlb(service *corev1.Service) *types.Nftlb {
	nftlb := &types.Nftlb{
		Farms: make([]types.Farm, len(service.Spec.Ports)),
	}

	// Read the annotations collected in the "annotations" field of the service
	annotations := getAnnotations(service)

	// Read useful values from the Service to be passed to servicePortAsAddress() instead of passing the Service
	serviceData := &types.ServiceData{
		Name:        service.Name,
		ClusterIP:   service.Spec.ClusterIP,
		Type:        string(service.Spec.Type),
		Family:      findFamily(service),
		ExternalIPs: service.Spec.ExternalIPs,
	}

	// Make wait group to syncronize every ServicePort
	wg := new(sync.WaitGroup)
	wg.Add(len(service.Spec.Ports))

	// Initialize a farm name slice based on the Service name (this is needed when a Service is deleted)
	farmsPerService[service.Name] = make([]string, len(service.Spec.Ports))

	// 1 ServicePort (k8s) = 1 Farm + 1 Address/Farm (nftlb)
	for index := range service.Spec.Ports {
		// Process all ServicePorts in parallel, using goroutines
		go func(servicePort *corev1.ServicePort, index int) {
			// Release lock after this func has finished
			defer wg.Done()

			// Parse ServicePort as Farm
			farm := servicePortAsFarm(servicePort, serviceData, annotations)

			// Set it in the Farms slice
			nftlb.Farms[index] = *farm

			// Branch out the non critical path (map assignments, DSR mode)
			go nonCriticalPathService(farm, service, index)
		}(&service.Spec.Ports[index], index)
	}

	// Wait until all locks are released
	wg.Wait()

	// Return a filled Nftlb struct
	return nftlb
}

// servicePortAsFarm returns a Farm struct filled with data from a ServicePort and some ServiceData values.
func servicePortAsFarm(servicePort *corev1.ServicePort, serviceData *types.ServiceData, annotations *types.Annotations) *types.Farm {
	address := types.Address{
		Family:   serviceData.Family,
		Protocol: strings.ToLower(string(servicePort.Protocol)),
	}

	if serviceData.Type == "ClusterIP" {
		// If the Service type is ClusterIP, add name, Service ClusterIP as ip-addr and ports
		address.Name = FormatName(serviceData.Name, servicePort.Name)
		address.IPAddr = serviceData.ClusterIP
		address.Ports = strconv.FormatInt(int64(servicePort.Port), 10)
	} else {
		// If the Service type is NodePort, add name and NodePort port ("ip-addr" is empty)
		address.Name = FormatNodePortName(serviceData.Name, servicePort.Name)
		address.Ports = strconv.FormatInt(int64(servicePort.NodePort), 10)
	}

	farm := &types.Farm{
		Name:         FormatName(serviceData.Name, servicePort.Name),
		Mode:         annotations.Mode,
		Persistence:  annotations.Persistence,
		PersistTTL:   annotations.PersistTTL,
		Scheduler:    annotations.Scheduler,
		SchedParam:   annotations.SchedParam,
		Helper:       annotations.Helper,
		Log:          annotations.Log,
		LogPrefix:    annotations.LogPrefix,
		EstConnlimit: annotations.EstConnlimit,
		Iface:        annotations.Iface,
		IntraConnect: "on",
		State:        "up",
		Addresses: []types.Address{
			address,
		},
	}

	// Add externalIPs as addresses
	for index, externalIP := range serviceData.ExternalIPs {
		farm.Addresses = append(farm.Addresses, types.Address{
			Family:   serviceData.Family,
			Protocol: strings.ToLower(string(servicePort.Protocol)),
			Name:     FormatExternalIPName(serviceData.Name, servicePort.Name, index),
			IPAddr:   externalIP,
			Ports:    strconv.FormatInt(int64(servicePort.Port), 10),
		})
	}

	return farm
}

func nonCriticalPathService(farm *types.Farm, service *corev1.Service, index int) {
	// Set farm name in map
	farmsPerService[service.Name][index] = farm.Name

	// Set addresses names in map
	addressesPerFarm[farm.Name] = make([]string, len(farm.Addresses))
	for idxAddress, address := range farm.Addresses {
		addressesPerFarm[farm.Name][idxAddress] = address.Name
	}

	// DSR mode
	if farm.Mode == "dsr" {
		// Enable DSR for future backends (for each backend, an interface is made)
		dsr.Enable(farm)
	} else if dsr.IsEnabled(farm) {
		// If this farm had DSR enabled before and no DSR annotation is specified, it means that DSR
		// must be disabled and interfaces must be deleted
		dsr.Disable(farm)
	}
}

func findFamily(service *corev1.Service) string {
	if localhostIP := net.ParseIP(service.Spec.ClusterIP); localhostIP.To4() != nil {
		return "ipv4"
	}
	return "ipv6"
}

func findIface(mode string) string {
	if mode == "dsr" {
		return config.DockerInterfaceBridge
	}
	return ""
}
