package json

import (
	"fmt"

	"github.com/zevenet/kube-nftlb/pkg/auth"
	"github.com/zevenet/kube-nftlb/pkg/dsr"
	"github.com/zevenet/kube-nftlb/pkg/types"

	corev1 "k8s.io/api/core/v1"
)

var (
	clienset = auth.GetClientset()

	// Active endpoints per service
	portsMap = map[string][]string{}

	// Active nodeports per service
	nodePortMap = map[string][]string{}

	// Active externalIPs per service
	externalIPsMap = map[string][]string{}

	// Active number of maxconns per service in the backend
	maxConnsMap = map[string]string{}
)

// ParseEndpointsAsFarms
func ParseEndpointsAsFarms(endpoints *corev1.Endpoints) types.Farms {
	endpointName := endpoints.ObjectMeta.Name
	state := "up"
	maxconns := "0"
	var farmsSlice []types.Farm
	// Go through each of the services we have created, specifically for each ip.
	// Then create and assign a backend for each port of our service.
	for _, endpoint := range endpoints.Subsets {
		backendID := 0
		for _, address := range endpoint.Addresses {
			// Get the ip and the name of the backends. If the name field is empty, it is assigned the same as the service.
			ipBackend := address.IP
			backendName := endpoints.ObjectMeta.Name
			if address.TargetRef != nil {
				backendName = address.TargetRef.Name
			}

			// We proceed to create each of the backends
			for _, port := range endpoint.Ports {
				// If the port name field is empty, it is assigned one by default
				if port.Name == "" {
					port.Name = "default"
				}

				// Attach the backends to the service
				farmName := assignFarmNameService(endpointName, port.Name)

				// Get backend ports
				portBackend := fmt.Sprint(port.Port)
				if dsr.ExistsServiceDSR(farmName) {
					dsr.AddInterfaceDsr(farmName, backendName, dsr.GetVirtualAddr(farmName))
					portBackend = dsr.GetVirtualPorts(farmName)
				}

				//
				if maxConnsMap[farmName] != "0" {
					maxconns = maxConnsMap[farmName]
				}

				backends := make([]types.Backend, 0)
				backend := createBackend(backendName, ipBackend, state, portBackend, maxconns)
				backends = append(backends, backend)
				var farm = types.Farm{
					Name:     fmt.Sprintf("%s", farmName),
					Backends: backends,
				}
				farmsSlice = append(farmsSlice, farm)
				// Check if the current service is of type nodePort thanks to the global variable that we have created previously
				// If this is the case, the nodePort service is assigned the same backends as the original service
				nodePortFarmName := assignFarmNameNodePort(farmName, "nodePort")
				if /*Contains(nodePortArray, nodePortFarmName)*/ true {
					var farm = types.Farm{
						Name:     fmt.Sprintf("%s", nodePortFarmName),
						Backends: backends,
					}
					farmsSlice = append(farmsSlice, farm)
				}

				// Check if the current service has externalIPs thanks to the global variable that we have created previously
				// If this is the case, the externalIPs service is assigned the same backends as the original service
				if _, ok := externalIPsMap[farmName]; ok {
					for _, farmExternalIPs := range externalIPsMap[farmName] {
						var farm = types.Farm{
							Name:     fmt.Sprintf("%s", farmExternalIPs),
							Backends: backends,
						}
						farmsSlice = append(farmsSlice, farm)
					}
				}

				backendID++
			}
		}
	}
	// Returns the filled struct
	return types.Farms{
		Farms: farmsSlice,
	}
}

// DeleteEndpointsBackends returns every backend name through a channel from a farm given a Endpoints object.
func DeleteEndpointsBackends(endpoints *corev1.Endpoints, backendPathsChan chan<- string) {
	farmName := endpoints.ObjectMeta.Name

	for _, backendName := range portsMap[farmName] {
		backendPathsChan <- fmt.Sprintf("farms/%s/backends/%s", farmName, backendName)
	}
	go delete(portsMap, farmName)

	for _, backendName := range nodePortMap[farmName] {
		backendPathsChan <- fmt.Sprintf("farms/%s/backends/%s", farmName, backendName)
	}
	go delete(nodePortMap, farmName)

	for _, backendName := range externalIPsMap[farmName] {
		backendPathsChan <- fmt.Sprintf("farms/%s/backends/%s", farmName, backendName)
	}
	go delete(externalIPsMap, farmName)

	close(backendPathsChan)
}

func GetEndpointsFarmName(endpoints *corev1.Endpoints) string {
	return endpoints.ObjectMeta.Name
}
