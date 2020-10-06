package parser

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/zevenet/kube-nftlb/pkg/config"
	"github.com/zevenet/kube-nftlb/pkg/types"

	corev1 "k8s.io/api/core/v1"
)

// ServiceAsPaths sends farm and addresses paths through a channel to the controller. The controller then sends a
// DELETE request to nftlb for every path.
func ServiceAsPaths(service *corev1.Service, pathChan chan<- string) {
	// Send farm and addresses paths to the controller
	pathChan <- fmt.Sprintf("farms/%s", service.Name)
	for _, addressName := range addressesPerService[service.Name] {
		pathChan <- fmt.Sprintf("addresses/%s", addressName)
	}

	// Remove from memory addresses names mapped to this Service
	delete(addressesPerService, service.Name)

	close(pathChan)
}

// ServiceAsNftlb analyzes a Service and returns a filled Nftlb struct.
func ServiceAsNftlb(service *corev1.Service) *types.Nftlb {
	// Read the annotations collected in the "annotations" field of the service
	annotations := getAnnotations(service)

	// Read useful values from the Service to be passed to servicePortAsAddress() instead of passing the Service
	serviceData := &types.ServiceData{
		Name:      service.Name,
		ClusterIP: service.Spec.ClusterIP,
		Type:      string(service.Spec.Type),
		Family:    findFamily(service),
	}

	// 1 Service (k8s) = 1 Farm (nftlb)
	farm := types.Farm{
		Name:         service.Name,
		Mode:         annotations.Mode,
		Persistence:  annotations.Persistence,
		PersistTTL:   annotations.PersistTTL,
		Scheduler:    annotations.Scheduler,
		SchedParam:   annotations.SchedParam,
		Helper:       annotations.Helper,
		Log:          annotations.Log,
		LogPrefix:    annotations.LogPrefix,
		EstConnlimit: annotations.EstConnlimit,
		Iface:        findIface(annotations.Mode),
		IntraConnect: "on",
		State:        "up",
		Addresses:    make([]types.Address, len(service.Spec.Ports)),
	}

	// Initialize an address name slice based on the Service name
	addressesPerService[service.Name] = make([]string, len(service.Spec.Ports))

	// Make wait group to syncronize every ServicePort
	wg := new(sync.WaitGroup)
	wg.Add(len(service.Spec.Ports))

	// 1 ServicePort (k8s) = 1 Address (nftlb)
	for index := range service.Spec.Ports {
		// Process all ServicePorts in parallel, using goroutines
		go func(servicePort *corev1.ServicePort, index int) {
			// Release lock after this func has finished
			defer wg.Done()

			// Parse ServicePort as Address
			address := servicePortAsAddress(servicePort, serviceData)

			// Append address to farm addresses
			farm.Addresses[index] = address

			// Append address name to addresses based on this Service
			addressesPerService[service.Name][index] = address.Name
		}(&service.Spec.Ports[index], index)
	}

	// Wait until all locks are released
	wg.Wait()

	// Return a filled Nftlb struct
	return &types.Nftlb{
		Farms: []types.Farm{
			farm,
		},
	}
}

// servicePortAsAddress returns an Address struct filled with data from a ServicePort and some ServiceData values.
func servicePortAsAddress(servicePort *corev1.ServicePort, serviceData *types.ServiceData) types.Address {
	// Make Address for this ServicePort
	address := types.Address{
		Family:   serviceData.Family,
		Protocol: strings.ToLower(string(servicePort.Protocol)),
	}

	if serviceData.Type == "ClusterIP" {
		// If the Service type is ClusterIP, add name, Service ClusterIP as ip-addr and ports
		address.Name = FormatName(serviceData.Name, servicePort.Name)
		address.IPAddr = serviceData.ClusterIP
		address.Ports = strconv.FormatInt(int64(servicePort.Port), 10)
	} else if serviceData.Type != "ExternalName" {
		// If the Service type isn't ExternalName, add name and NodePort port ("ip-addr" is empty)
		address.Name = FormatNodePortName(serviceData.Name, servicePort.Name)
		address.Ports = strconv.FormatInt(int64(servicePort.NodePort), 10)
	}

	return address
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
