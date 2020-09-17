package parser

import (
	"fmt"
	"sync"

	"github.com/zevenet/kube-nftlb/pkg/dsr"
	"github.com/zevenet/kube-nftlb/pkg/types"

	corev1 "k8s.io/api/core/v1"
)

// EndpointsAsFarms reads a Endpoints and returns a filled Farms struct.
func EndpointsAsFarms(endpoints *corev1.Endpoints) *types.Farms {
	farms := make([]types.Farm, 0)
	var wg sync.WaitGroup

	for idxSubset, subset := range endpoints.Subsets {
		for idxAddress := range subset.Addresses {
			// Add a lock for every EndpointAddress
			wg.Add(1)
			go CreateEndpointAddress(endpoints, &subset.Addresses[idxAddress], &endpoints.Subsets[idxSubset], &farms, &wg)
		}
	}

	// Wait until all locks are released
	wg.Wait()

	// Returns the filled struct
	return &types.Farms{
		Farms: farms,
	}
}

// CreateEndpointAddress appends farms parsed from a EndpointAddress to a Farm slice.
func CreateEndpointAddress(endpoints *corev1.Endpoints, address *corev1.EndpointAddress, subset *corev1.EndpointSubset, farms *[]types.Farm, wg *sync.WaitGroup) {
	// If the name field is empty, it is assigned the same as the service.
	backendName := endpoints.Name
	if address.TargetRef != nil {
		backendName = address.TargetRef.Name
	}

	// We proceed to create each of the backends
	for _, endpointPort := range subset.Ports {
		farm := types.Farm{
			Name: FormatFarmName(endpoints.Name, endpointPort.Name),
		}

		// Get backend ports
		portBackend := fmt.Sprint(endpointPort.Port)
		if dsr.ExistsServiceDSR(farm.Name) {
			portBackend = dsr.GetVirtualPorts(farm.Name)
			go dsr.AddInterfaceDsr(farm.Name, backendName, dsr.GetVirtualAddr(farm.Name))
		}

		farm.Backends = []types.Backend{
			{
				Name:         backendName,
				IPAddr:       address.IP,
				State:        "up",
				Port:         portBackend,
				EstConnlimit: maxConnsMap[farm.Name],
			},
		}
		*farms = append(*farms, farm)
		go func() {
			farmsPerEndpointMap[endpoints.Name] = append(farmsPerEndpointMap[endpoints.Name], farm.Name)
			backendsPerFarm[farm.Name] = append(backendsPerFarm[farm.Name], backendName)
		}()

		// The NodePort Service is assigned the same backends as the original service
		if farmNameNodePort := FormatNodePortFarmName(endpoints.Name, endpointPort.Name); existsFarm(endpoints.Name, farmNameNodePort) {
			farm.Name = farmNameNodePort
			*farms = append(*farms, farm)
			go func() {
				farmsPerEndpointMap[endpoints.Name] = append(farmsPerEndpointMap[endpoints.Name], farm.Name)
				backendsPerFarm[farm.Name] = append(backendsPerFarm[farm.Name], backendName)
			}()
		}

		// The ExternalIP Service is assigned the same backends as the original service
		for _, farmNameExternalIP := range farmsPerExternalIPResourceMap[endpoints.Name] {
			farm.Name = farmNameExternalIP
			*farms = append(*farms, farm)
			go func() {
				farmsPerEndpointMap[endpoints.Name] = append(farmsPerEndpointMap[endpoints.Name], farm.Name)
				backendsPerFarm[farm.Name] = append(backendsPerFarm[farm.Name], backendName)
			}()
		}
	}

	// Release lock
	wg.Done()
}

// DeleteEndpointsBackends sends backend paths through a channel to the controller. The controller then deletes every backend.
func DeleteEndpointsBackends(endpoints *corev1.Endpoints, backendPathsChan chan<- string) {
	for _, farmName := range farmsPerEndpointMap[endpoints.Name] {
		for _, backendName := range backendsPerFarm[farmName] {
			backendPathsChan <- fmt.Sprintf("farms/%s/backends/%s", farmName, backendName)
		}
		backendsPerFarm[farmName] = []string{}
	}
	farmsPerEndpointMap[endpoints.Name] = []string{}
	close(backendPathsChan)
}
