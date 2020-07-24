package json

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	configFarm "github.com/zevenet/kube-nftlb/pkg/farms"
	types "github.com/zevenet/kube-nftlb/pkg/types"
	v1 "k8s.io/api/core/v1"
)

// Check if the service has active nodeports. If that's the case, store it in the list.
var nodePortArray []string

// EncodeJSON returns a JSONnftlb struct with its fields filled with the JSON data.
func EncodeJSON(stringJSON string) types.JSONnftlb {
	var encodedJSON types.JSONnftlb
	if err := json.Unmarshal([]byte(stringJSON), &encodedJSON); err != nil {
		panic(err.Error())
	}
	return encodedJSON
}

// GetJSONnftlbFromService returns a JSONnftlb struct filled with any Service data.
func GetJSONnftlbFromService(service *v1.Service) types.JSONnftlb {
	// Gets persistence and Stickiness timeout in seconds
	persistence, persistenceTTL := getPersistence(service)
	// Read the annotations collected in the "annotations" field of the service
	mode, scheduler, helper, log, logprefix := getAnnotations(service)
	// Creates the service
	var farmsSlice []types.Farm
	serviceName := service.ObjectMeta.Name
	farmName := ""
	// Default values of service (Family, State and Intraconnect)
	family := "ipv4"
	state := "up"
	intraconnect := "on"
	// When creating services we can create several from the same yaml configuration file
	// For this we take into account the port field of the yaml configuration file. We create a service for each name field in ports
	for _, port := range service.Spec.Ports {
		// // Gets the name, protocol, port and ip of the service. If we are creating a single service and it does not have a port name, we assign it a default name
		if port.Name == "" {
			farmName = configFarm.AssignFarmNameService(serviceName, "default")
		} else {
			farmName = configFarm.AssignFarmNameService(serviceName, port.Name)
		}
		nameProtocol := strings.ToLower(string(port.Protocol))
		portString := fmt.Sprint(port.Port)
		virtualAddr := service.Spec.ClusterIP
		// Creates and fill the farm.
		var farm = createFarm(farmName, family, virtualAddr, portString, mode, nameProtocol, scheduler, helper, log, logprefix, state, intraconnect, persistence, persistenceTTL, types.Backends{})
		farmsSlice = append(farmsSlice, farm)
		// Check if the service is type NodePort
		// If so, modify the name of the original service, modify the port, modify its virtualip and store its name in a global variable to then be able to reference it
		if service.Spec.Type == "NodePort" || service.Spec.Type == "LoadBalancer" && port.NodePort >= 0 {
			farmName = configFarm.AssignFarmNameNodePort(serviceName+"--"+port.Name, "nodePort")
			virtualAddr = ""
			portString = fmt.Sprint(port.NodePort)
			// Creates and fills the NodePort farm
			var farm = createFarm(farmName, family, virtualAddr, portString, mode, nameProtocol, scheduler, helper, log, logprefix, state, intraconnect, persistence, persistenceTTL, types.Backends{})
			farmsSlice = append(farmsSlice, farm)
			nodePortArray = append(nodePortArray, farmName)
		}
	}
	// Returns the filled struct
	return types.JSONnftlb{
		Farms: farmsSlice,
	}
}

// GetJSONnftlbFromEndpoints returns a JSONnftlb struct filled with any Endpoints data.
func GetJSONnftlbFromEndpoints(endpoints *v1.Endpoints) types.JSONnftlb {
	objName := endpoints.ObjectMeta.Name
	portName := ""
	farmName := ""
	state := "up"
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
			if address.TargetRef != nil {
				backendName = address.TargetRef.Name
			} else if address.TargetRef == nil {
				backendName = endpoints.ObjectMeta.Name
			}
			// We proceed to create each of the backends
			for _, port := range endpoint.Ports {
				portBackend := fmt.Sprint(port.Port)
				var backends types.Backends
				var backend = createBackend(backendName, ipBackend, state, portBackend)
				// If the port name field is empty, it is assigned one by default.
				// Once done, attach the backends to the service.
				backends = append(backends, backend)
				if port.Name == "" {
					portName = "default"
				} else {
					portName = port.Name
				}
				farmName = configFarm.AssignFarmNameService(objName, portName)
				var farm = types.Farm{
					Name:     fmt.Sprintf("%s", farmName),
					Backends: backends,
				}
				farmsSlice = append(farmsSlice, farm)
				// Check if the current service is of type nodePort thanks to the global variable that we have created previously
				// If this is the case, the nodePort service is assigned the same backends as the original service
				farmName = configFarm.AssignFarmNameNodePort(farmName, "nodePort")
				if Contains(nodePortArray, farmName) {
					var farm = types.Farm{
						Name:     fmt.Sprintf("%s", farmName),
						Backends: backends,
					}
					farmsSlice = append(farmsSlice, farm)
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

func createFarm(farmName string, family string, virtualAddr string, virtualPorts string, mode string, protocol string, scheduler string, helper string, log string, logPrefix string, state string, intraconnect string, persistence string, persistTTL string, backends types.Backends) types.Farm {
	// we create the farm based on the types that we have previously defined in types/json
	var farmCreated = types.Farm{
		Name:         farmName,
		Family:       family,
		VirtualAddr:  fmt.Sprintf("%s", virtualAddr),
		VirtualPorts: fmt.Sprintf("%s", virtualPorts),
		Mode:         mode,
		Protocol:     protocol,
		Scheduler:    scheduler,
		Helper:       helper,
		Log:          log,
		LogPrefix:    logPrefix,
		State:        state,
		Intraconnect: intraconnect,
		Persistence:  fmt.Sprintf("%s", persistence),
		PersistTTL:   fmt.Sprintf("%s", persistTTL),
		Backends:     backends,
	}
	return farmCreated
}

func createBackend(name string, ipAddr string, state string, port string) types.Backend {
	var backendCreated = types.Backend{
		Name:   fmt.Sprintf("%s", name),
		IPAddr: ipAddr,
		State:  state,
		Port:   port,
	}
	return backendCreated
}

func getPersistence(service *v1.Service) (string, string) {
	persistence := ""
	persistenceTTL := ""
	if service.Spec.SessionAffinity == "ClientIP" {
		persistence = "srcip"
	} else if service.Spec.SessionAffinity == "None" {
		persistence = "none"
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

func getAnnotations(service *v1.Service) (string, string, string, string, string) {
	// First try reading the annotations for fields that can be configured in the nftlb service
	// If there are no annotations for all the fields, default values ​​are set.
	// You don't need to worry about sending empty variables as it is configured so if a variable is sent empty it is not included in the json that configures the nftlb service.
	mode := "snat"
	scheduler := ""
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
				mode = value
			} else if strings.ToLower(string(field[0])) == "scheduler" {
				scheduler = value
			} else if strings.ToLower(string(field[0])) == "helper" {
				helper = value
			} else if strings.ToLower(string(field[0])) == "log" {
				log = value
			} else if strings.ToLower(string(field[0])) == "logprefix" && log != "" && log != "none" {
				logprefix = value
			}
		}
	}
	return mode, scheduler, helper, log, logprefix
}

func GetNodePortArray() []string {
	return nodePortArray
}
