package parser

import (
	"fmt"
	"sync"

	"github.com/zevenet/kube-nftlb/pkg/types"

	corev1 "k8s.io/api/core/v1"
)

// EndpointsAsPaths sends backend paths through a channel to the controller. The controller then deletes every backend.
func EndpointsAsPaths(endpoints *corev1.Endpoints, pathsChan chan<- string) {
	for _, backendName := range backendsPerEndpoints[endpoints.Name] {
		pathsChan <- fmt.Sprintf("farms/%s/backends/%s", endpoints.Name, backendName)
	}

	delete(backendsPerEndpoints, endpoints.Name)
	close(pathsChan)
}

// EndpointsAsNftlb reads a Endpoints object and returns a filled Nftlb struct.
func EndpointsAsNftlb(endpoints *corev1.Endpoints) *types.Nftlb {
	farm := types.Farm{
		Name:     endpoints.Name,
		Backends: make([]types.Backend, 0),
	}
	backendsPerEndpoints[endpoints.Name] = make([]string, 0)

	var wg sync.WaitGroup

	for idxSubset, subset := range endpoints.Subsets {
		for idxPort := range subset.Ports {
			for idxAddress := range subset.Addresses {
				// Add a lock for every EndpointAddress
				wg.Add(1)

				// 1 EndpointAddress for every EndpointPort in a EndpointSubset (k8s) = 1 Backend (nftlb)
				go func(epSubset *corev1.EndpointSubset, epPort *corev1.EndpointPort, epAddress *corev1.EndpointAddress, idxSubset int, idxPort int, idxAddress int) {
					defer wg.Done()

					backend := types.Backend{
						Name:   fmt.Sprintf("%s--subset-%d--port-%d--address-%d", endpoints.Name, idxSubset, idxPort, idxAddress),
						IPAddr: epAddress.IP,
						State:  "up",
						Port:   fmt.Sprint(epPort.Port),
					}

					farm.Backends = append(farm.Backends, backend)

					backendsPerEndpoints[endpoints.Name] = append(backendsPerEndpoints[endpoints.Name], backend.Name)
				}(&endpoints.Subsets[idxSubset], &subset.Ports[idxPort], &subset.Addresses[idxAddress], idxSubset, idxPort, idxAddress)
			}
		}
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
