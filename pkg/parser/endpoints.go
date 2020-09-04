package parser

import (
	"context"
	"fmt"

	"github.com/zevenet/kube-nftlb/pkg/dsr"
	"github.com/zevenet/kube-nftlb/pkg/types"
	"k8s.io/apimachinery/pkg/labels"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EndpointsAsFarms
func EndpointsAsFarms(endpoints *corev1.Endpoints) *types.Farms {
	endpointName := endpoints.Name
	serviceList := getServiceListFromEndpoints(endpoints.Labels)
	state := "up"
	maxconns := "0"
	farms := make([]types.Farm, 0)

	for _, service := range serviceList.Items {
		serviceName := service.Name
		for _, endpoint := range endpoints.Subsets {
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
					if maxConnsMap[serviceName][farmName] != "0" {
						maxconns = maxConnsMap[serviceName][farmName]
					}

					backends := make([]types.Backend, 0)
					backend := types.Backend{
						Name:         backendName,
						IPAddr:       ipBackend,
						State:        state,
						Port:         portBackend,
						EstConnlimit: maxconns,
					}
					backends = append(backends, backend)
					var farm = types.Farm{
						Name:     farmName,
						Backends: backends,
					}
					farms = append(farms, farm)

					// Check if the current service is of type nodePort thanks to the global variable that we have created previously
					// If this is the case, the nodePort service is assigned the same backends as the original service
					nodePortFarmName := assignFarmNameNodePort(farmName, "nodePort")
					if /*Contains(nodePortArray, nodePortFarmName)*/ true {
						var farm = types.Farm{
							Name:     nodePortFarmName,
							Backends: backends,
						}
						farms = append(farms, farm)
					}

					// Check if the current service has externalIPs thanks to the global variable that we have created previously
					// If this is the case, the externalIPs service is assigned the same backends as the original service
					if _, ok := externalIPsMap[serviceName]; ok {
						for _, farmNameExternalIP := range externalIPsMap[farmName] {
							var farm = types.Farm{
								Name:     farmNameExternalIP,
								Backends: backends,
							}
							farms = append(farms, farm)
						}
					}
				}
			}
		}
	}

	// Returns the filled struct
	return &types.Farms{
		Farms: farms,
	}
}

// DeleteEndpointsBackends returns every backend name through a channel from a farm given a Endpoints object.
func DeleteEndpointsBackends(endpoints *corev1.Endpoints, backendPathsChan chan<- string) {
	serviceList := getServiceListFromEndpoints(endpoints.Labels)

	for _, service := range serviceList.Items {
		farmName := service.Name

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
	}

	close(backendPathsChan)
}

func getServiceListFromEndpoints(endpointsLabels map[string]string) *corev1.ServiceList {
	opts := metav1.ListOptions{
		LabelSelector: labels.Set(endpointsLabels).String(),
	}

	list, _ := clientset.CoreV1().Services(corev1.NamespaceAll).List(context.TODO(), opts)

	return list
}
