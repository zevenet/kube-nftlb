package parser

import (
	"fmt"
	"sync"

	"github.com/zevenet/kube-nftlb/pkg/types"

	corev1 "k8s.io/api/core/v1"
)

// EndpointsAsPaths sends backend paths through a channel to the controller. The controller then deletes every backend.
func EndpointsAsPaths(endpoints *corev1.Endpoints, pathsChan chan<- string) {
	for _, farmName := range farmsPerService[endpoints.Name] {
		for _, backendName := range backendsPerFarm[farmName] {
			pathsChan <- fmt.Sprintf("farms/%s/backends/%s", farmName, backendName)
		}
		delete(backendsPerFarm, farmName)
	}

	close(pathsChan)
}

// EndpointsAsNftlb reads a Endpoints object and returns a filled Nftlb struct.
func EndpointsAsNftlb(endpoints *corev1.Endpoints) *types.Nftlb {
	nftlb := &types.Nftlb{
		Farms: make([]types.Farm, 0),
	}

	for idxSubset, subset := range endpoints.Subsets {
		// 1 EndpointPort (k8s) = 1 Farm (nftlb)
		for idxPort, port := range subset.Ports {
			farm := types.Farm{
				Name:     FormatName(endpoints.Name, port.Name),
				Backends: make([]types.Backend, len(subset.Addresses)),
			}
			backendsPerFarm[farm.Name] = make([]string, len(subset.Addresses))

			// Add a lock for every EndpointAddress
			wg := new(sync.WaitGroup)
			wg.Add(len(subset.Addresses))

			// 1 EndpointAddress for every EndpointPort in a EndpointSubset (k8s) = 1 Backend (nftlb)
			for idxAddress := range subset.Addresses {
				go func(epPort *corev1.EndpointPort, epAddress *corev1.EndpointAddress, idxAddress int) {
					defer wg.Done()

					backend := types.Backend{
						IPAddr: epAddress.IP,
						State:  "up",
						Port:   fmt.Sprint(epPort.Port),
					}

					if epAddress.TargetRef != nil {
						backend.Name = FormatName(epAddress.TargetRef.Name, epPort.Name)
					} else {
						backend.Name = FormatName(endpoints.Name, epPort.Name)
					}

					farm.Backends[idxAddress] = backend
					backendsPerFarm[farm.Name][idxAddress] = backend.Name
				}(&endpoints.Subsets[idxSubset].Ports[idxPort], &endpoints.Subsets[idxSubset].Addresses[idxAddress], idxAddress)
			}

			// Wait until all EndpointPort locks are released
			wg.Wait()

			nftlb.Farms = append(nftlb.Farms, farm)
		}
	}

	// Return a filled Nftlb struct
	return nftlb
}
