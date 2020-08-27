package json

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strings"
	"strconv"

	config "github.com/zevenet/kube-nftlb/pkg/config"
	dsr "github.com/zevenet/kube-nftlb/pkg/dsr"
	configFarm "github.com/zevenet/kube-nftlb/pkg/farms"
	types "github.com/zevenet/kube-nftlb/pkg/types"
	v1 "k8s.io/api/core/v1"
	kubernetes "k8s.io/client-go/kubernetes"
)

// Check if the service has active nodeports. If that's the case, store it in the list.
var nodePortArray []string

// Check if the service has type DSR. If that's the case, store in the list
var serviceDsr = make(map[string]*types.DsrStruct)

// Check if the service has active one number of maxconns in the backend
var maxConnsMap = map[string]string{}

// Check if the service has active externalIPs
var externalIPsMap = map[string][]string{}

// EncodeJSON returns a JSONnftlb struct with its fields filled with the JSON data.
func EncodeJSON(stringJSON string) types.JSONnftlb {
	var encodedJSON types.JSONnftlb
	if err := json.Unmarshal([]byte(stringJSON), &encodedJSON); err != nil {
		panic(err.Error())
	}
	return encodedJSON
}

// GetJSONnftlbFromService returns a JSONnftlb struct filled with any Service data.
func GetJSONnftlbFromService(service *v1.Service, clientset *kubernetes.Clientset) types.JSONnftlb {
	// Gets persistence and Stickiness timeout in seconds
	persistence, persistenceTTL := getPersistence(service)
	findMaxConns(service)
	// Creates the service
	var farmsSlice []types.Farm
	serviceName := service.ObjectMeta.Name
	farmName := ""
	// Default values of service (State, Intraconnect and Iface)
	state := "up"
	intraconnect := "on"
	iface := ""
	// find out the ip family, by default the values ​​is ipv4
	family := findFamily(service)
	// When creating services we can create several from the same yaml configuration file
	// For this we take into account the port field of the yaml configuration file. We create a service for each name field in ports
	for _, port := range service.Spec.Ports {
		// // Gets the name, protocol, port and ip of the service. If we are creating a single service and it does not have a port name, we assign it a default name
		if port.Name == "" {
			farmName = configFarm.AssignFarmNameService(serviceName, "default")
		} else {
			farmName = configFarm.AssignFarmNameService(serviceName, port.Name)
		}
		// Read the annotations collected in the "annotations" field of the service
		mode, scheduler, schedulerParam, helper, log, logprefix := getAnnotations(service, farmName)
		nameProtocol := strings.ToLower(string(port.Protocol))
		portString := fmt.Sprint(port.Port)
		virtualAddr := service.Spec.ClusterIP
		// check if the service has mode dsr, if that's the case, append to the list
		if mode == "dsr" {
			if _, ok := serviceDsr[farmName]; !ok {
				serviceDsr[farmName] = &types.DsrStruct{}
				serviceDsr = dsr.CreateInterfaceDsr(farmName, service, clientset, serviceDsr)
			}
			serviceDsr[farmName].VirtualAddr = virtualAddr
			serviceDsr[farmName].VirtualPorts = portString
			iface = config.DockerInterfaceBridge
		} else if mode != "dsr" {
			if _, ok := serviceDsr[farmName]; ok {
				if service.Spec.Type != "NodePort" {
					serviceDsr = dsr.DeleteInterfaceDsr(farmName, serviceDsr)
				}
			}
		}
		// Creates and fill the farm.
		var farm = createFarm(farmName, family, virtualAddr, portString, mode, nameProtocol, scheduler, schedulerParam, helper, log, logprefix, state, intraconnect, persistence, persistenceTTL, iface, types.Backends{})
		farmsSlice = append(farmsSlice, farm)
		// Check if the service is type NodePort
		// If so, modify the name of the original service, modify the port, modify its virtualip and store its name in a global variable to then be able to reference it
		if service.Spec.Type == "NodePort" || service.Spec.Type == "LoadBalancer" && port.NodePort >= 0 {
			farmName = configFarm.AssignFarmNameNodePort(serviceName+"--"+port.Name, "nodePort")
			virtualAddr = ""
			portString = fmt.Sprint(port.NodePort)
			// check if the service has mode dsr, if that's the case, append to the list
			if mode == "dsr" {
				if _, ok := serviceDsr[farmName]; !ok {
					serviceDsr[farmName] = &types.DsrStruct{}
					serviceDsr = dsr.CreateInterfaceDsr(farmName, service, clientset, serviceDsr)
				}
				serviceDsr[farmName].VirtualAddr = virtualAddr
				serviceDsr[farmName].VirtualPorts = portString
				iface = config.DockerInterfaceBridge
			} else if mode != "dsr" {
				if _, ok := serviceDsr[farmName]; ok {
					if service.Spec.Type != "NodePort" {
						serviceDsr = dsr.DeleteInterfaceDsr(farmName, serviceDsr)
					}
				}
			}
			// Creates and fills the NodePort farm
			var farm = createFarm(farmName, family, virtualAddr, portString, mode, nameProtocol, scheduler, schedulerParam, helper, log, logprefix, state, intraconnect, persistence, persistenceTTL, iface, types.Backends{})
			farmsSlice = append(farmsSlice, farm)
			nodePortArray = append(nodePortArray, farmName)
		}

		// Check if the service has externalIPs field
		// We have to create an externalIPs service for each external ip that has been configured in the configuration file.yaml
		if len(service.Spec.ExternalIPs) >= 1 {
			for position, externalIPs := range service.Spec.ExternalIPs {
				virtualAddr = externalIPs
				externalIPsName := "externalIPs"+strconv.Itoa(position+1)
				externalIpFarmName := configFarm.AssignFarmNameExternalIPs(serviceName+"--"+port.Name, externalIPsName)
				var farm = createFarm(externalIpFarmName, family, virtualAddr, portString, mode, nameProtocol, scheduler, schedulerParam, helper, log, logprefix, state, intraconnect, persistence, persistenceTTL, iface, types.Backends{})
				farmsSlice = append(farmsSlice, farm)
				externalIPsMap[farmName] = append(externalIPsMap[farmName], externalIpFarmName)
			}
		}
	}
	// Returns the filled struct
	return types.JSONnftlb{
		Farms: farmsSlice,
	}
}

// GetJSONnftlbFromEndpoints returns a JSONnftlb struct filled with any Endpoints data.
func GetJSONnftlbFromEndpoints(endpoints *v1.Endpoints, clientset *kubernetes.Clientset) types.JSONnftlb {
	objName := endpoints.ObjectMeta.Name
	portName := ""
	farmName := ""
	state := "up"
	maxconns := "0"
	// Initializes farm/backends ID
	CreateFarmID(objName)
	var farmsSlice []types.Farm
	// Go through each of the services we have created, specifically for each ip.
	// Then create and assign a backend for each port of our service.
	for _, endpoint := range endpoints.Subsets {
		for _, address := range endpoint.Addresses {
			// Get the ip and the name of the backends. If the name field is empty, it is assigned the same as the service.
			ipBackend := address.IP
			backendName := ""
			//Uid := ""
			if address.TargetRef != nil {
				backendName = address.TargetRef.Name
			} else if address.TargetRef == nil {
				backendName = endpoints.ObjectMeta.Name
			}
			// We proceed to create each of the backends
			for _, port := range endpoint.Ports {
				// If the port name field is empty, it is assigned one by default.
				// Once done, attach the backends to the service.
				if port.Name == "" {
					portName = "default"
				} else {
					portName = port.Name
				}
				farmName = configFarm.AssignFarmNameService(objName, portName)
				portBackend := fmt.Sprint(port.Port)
				if maxConnsMap[farmName] != "0" {
					maxconns = maxConnsMap[farmName]
				}
				// check if the current service has mode dsr thanks to the global variable that we have created previously
				//If this is the case, with the dsr mode activated, the IP of the service is assigned to the loopback network interface of the backends
				if _, ok := serviceDsr[farmName]; ok {
					serviceDsr = dsr.AddInterfaceDsr(clientset, farmName, backendName, serviceDsr[farmName].VirtualAddr, serviceDsr)
					temp := serviceDsr[farmName].VirtualPorts
					portBackend = temp
				}
				var backends types.Backends
				var backend = createBackend(backendName, ipBackend, state, portBackend, maxconns)
				backends = append(backends, backend)
				var farm = types.Farm{
					Name:     fmt.Sprintf("%s", farmName),
					Backends: backends,
				}
				farmsSlice = append(farmsSlice, farm)
				if Contains(nodePortArray, farmName) {
					var farm = types.Farm{
						Name:     fmt.Sprintf("%s", farmName),
						Backends: backends,
					}
					farmsSlice = append(farmsSlice, farm)
				}
				// Check if the current service is of type nodePort thanks to the global variable that we have created previously
				// If this is the case, the nodePort service is assigned the same backends as the original service
				nodePortFarmName := configFarm.AssignFarmNameNodePort(farmName, "nodePort")
				if Contains(nodePortArray, nodePortFarmName) {
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
			}
		}
	}
	// Returns the filled struct
	return types.JSONnftlb{
		Farms: farmsSlice,
	}
}

// Contains returns true when "str" string is in "sl" slice.
func Contains(sl []string, str string) bool {
	for _, value := range sl {
		if value == str {
			return true
		}
	}
	return false
}

func createFarm(farmName string, family string, virtualAddr string, virtualPorts string, mode string, protocol string, scheduler string, schedulerParam string, helper string, log string, logPrefix string, state string, intraconnect string, persistence string, persistTTL string, iface string, backends types.Backends) types.Farm {
	// we create the farm based on the types that we have previously defined in types/json
	var farmCreated = types.Farm{
		Name:           farmName,
		Family:         family,
		Iface:          iface,
		VirtualAddr:    fmt.Sprintf("%s", virtualAddr),
		VirtualPorts:   fmt.Sprintf("%s", virtualPorts),
		Mode:           mode,
		Protocol:       protocol,
		Scheduler:      scheduler,
		SchedulerParam: schedulerParam,
		Helper:         helper,
		Log:            log,
		LogPrefix:      logPrefix,
		State:          state,
		Intraconnect:   intraconnect,
		Persistence:    fmt.Sprintf("%s", persistence),
		PersistTTL:     fmt.Sprintf("%s", persistTTL),
		Backends:       backends,
	}
	return farmCreated
}

func createBackend(name string, ipAddr string, state string, port string, maxconns string) types.Backend {
	var backendCreated = types.Backend{
		Name:            fmt.Sprintf("%s", name),
		IPAddr:          ipAddr,
		State:           state,
		Port:            port,
		BackendMaxConns: maxconns,
	}
	return backendCreated
}

func getPersistence(service *v1.Service) (string, string) {
	// First we get the persistence of our service. By default, annotations have priority ahead of the sessionAffinity and sessionAffinityConfig field.
	// If there are no annotations, the information in the sessionAffinity and sessionAffinityConfig field is collected.
	persistence := ""
	persistenceTTL := "60"
	var rgx = regexp.MustCompile(`[a-z]+$`)
	if service.ObjectMeta.Annotations != nil {
		for key, value := range service.ObjectMeta.Annotations {
			field := rgx.FindStringSubmatch(key)
			if strings.ToLower(string(field[0])) == "persistence" {
				//  check for multiple combination of values. Ej srcip+srcport && srcip+dstport
				splitField := strings.Split(value, "-")
				if len(splitField) > 1 {
					if splitField[0] == "srcip" && splitField[1] == "srcport" {
						persistence = splitField[0] + " " + splitField[1]
					} else if splitField[0] == "srcip" && splitField[1] == "dstport" {
						persistence = splitField[0] + " " + splitField[1]
					}
				} else if value == "srcip" || value == "dstip" || value == "srcport" || value == "dstport" || value == "srcmac" || value == "dstmac" {
					persistence = value
				}
			}
		}
	}
	if persistence == "" {
		if service.Spec.SessionAffinity == "ClientIP" {
			persistence = "srcip"
		} else if service.Spec.SessionAffinity == "None" {
			persistence = "none"
		}
	}
	if service.Spec.SessionAffinityConfig != nil {
		if service.Spec.SessionAffinityConfig.ClientIP != nil {
			if service.Spec.SessionAffinityConfig.ClientIP.TimeoutSeconds != nil {
				// Value between 0 and 86400 seconds (1 day max)
				persistenceTTL = fmt.Sprint(*(service.Spec.SessionAffinityConfig.ClientIP.TimeoutSeconds))
			}
		}
	}
	return persistence, persistenceTTL
}

func getAnnotations(service *v1.Service, farmName string) (string, string, string, string, string, string) {
	// First try reading the annotations for fields that can be configured in the nftlb service
	// If there are no annotations for all the fields, default values ​​are set.
	// You don't need to worry about sending empty variables as it is configured so if a variable is sent empty it is not included in the json that configures the nftlb service.
	mode := "snat"
	scheduler := "rr"
	sched_param := "none"
	helper := ""
	log := ""
	logprefix := ""
	// We use a regular expression to filter the string and get the field to configure in the annotations
	// Always respecting the format of the string | service.kubernetes.io/kube-nftlb-load-balancer-X | where X is the field to configure
	var rgx = regexp.MustCompile(`[a-z]+$`)
	if service.ObjectMeta.Annotations != nil {
		for key, value := range service.ObjectMeta.Annotations {
			field := rgx.FindStringSubmatch(key)
			if strings.ToLower(string(field[0])) == "mode" {
				if value == "snat" || value == "dnat" || value == "dsr" || value == "stlsdnat" || value == "local" {
					mode = value
				}
			} else if strings.ToLower(string(field[0])) == "scheduler" {
				rgx = regexp.MustCompile(`^[a-z]+`)
				field = rgx.FindStringSubmatch(value)
				if value == "rr" {
					scheduler = value
				} else if value == "symhash" {
					scheduler = value
				} else if strings.ToLower(string(field[0])) == "hash" {
					splitField := strings.Split(value, "-")
					// check for multiple combination of values. Ej srcip+srcport
					if len(splitField) > 2 {
						if splitField[1] == "srcip" && splitField[2] == "srcport" {
							sched_param = splitField[1] + " " + splitField[2]
						}
					} else {
						valueHash := splitField[1]
						if valueHash == "srcip" || valueHash == "dstip" || valueHash == "srcport" || valueHash == "dstport" || valueHash == "srcmac" || valueHash == "dstmac" {
							sched_param = valueHash
						}
						scheduler = "hash"
					}

				}
			} else if strings.ToLower(string(field[0])) == "helper" {
				if value == "amanda" || value == "ftp" || value == "h323" || value == "irc" || value == "netbios-ns" || value == "pptp" || value == "sane" || value == "sip" || value == "snmp" || value == "tftp" {
					helper = value
				}
			} else if strings.ToLower(string(field[0])) == "log" {
				if value == "none" || value == "forward" || value == "output" {
					log = value
					logprefix = farmName
				}
			} 
		}
	}
	return mode, scheduler, sched_param, helper, log, logprefix
}

func findMaxConns(service *v1.Service) {
	var farmSlice []string
	backendMaxConnsMap := "0"
	serviceName := service.ObjectMeta.Name
	var rgx = regexp.MustCompile(`[a-z]+$`)
	if service.ObjectMeta.Annotations != nil {
		farmName := ""
		for _, port := range service.Spec.Ports {
			if port.Name == "" {
				farmName = configFarm.AssignFarmNameService(serviceName, "default")
			} else {
				farmName = configFarm.AssignFarmNameService(serviceName, port.Name)
			}
			farmSlice = append(farmSlice, farmName)
		}
		for key, value := range service.ObjectMeta.Annotations {
			field := rgx.FindStringSubmatch(key)
			if strings.ToLower(string(field[0])) == "maxconns" {
				backendMaxConnsMap = value
			}
		}
	}
	for _, nameFarm := range farmSlice {
		maxConnsMap[nameFarm] = backendMaxConnsMap
	}
}

func GetNodePortArray() []string {
	// Returns the list of services that have nodeport. It is used to identify the nodeport services and then delete them together with the main service.
	return nodePortArray
}

func GetExternalIPsArray() map[string][]string {
	// Returns the list of services that have externalIPs. It is used to identify the externalIPs services and then delete them together with the main service.
	return externalIPsMap
}

func DeleteMaxConnsMap() {
	maxConnsMap = map[string]string{}
}

func findFamily(service *v1.Service) string {
	// Find out what type of version the service IP has, by default the value ​​is ipv4
	family := "ipv4"
	localhostIp := net.ParseIP(service.Spec.ClusterIP)
	if localhostIp.To4() != nil {
		family = "ipv4"
	} else if localhostIp.To16() != nil {
		family = "ipv6"
	}
	return family
}