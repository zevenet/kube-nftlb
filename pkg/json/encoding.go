package json

import (
	"encoding/json"
	"fmt"
	"strings"
	"regexp"

	types "github.com/zevenet/kube-nftlb/pkg/types"
	v1 "k8s.io/api/core/v1"
)

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
	// Extracts ports and protocols as strings
	protocolsSlice := map[string]int{"TCP": 0, "UDP": 0}
	var portsSlice []string
	// For every service port:
	for _, port := range service.Spec.Ports {
		// Gets the port as string
		portString := fmt.Sprint(port.Port)
		// Increase number of protocol ocurrences
		protocolsSlice[string(port.Protocol)]++
		// Don't duplicate if it exists already
		if !Contains(portsSlice, portString) {
			portsSlice = append(portsSlice, portString)
		}
	}
	// Ports slice -> port(s) string
	ports := strings.Join(portsSlice, ", ")
	// Gets the protocol
	protocol := ""
	if protocolsSlice["TCP"] > 0 && protocolsSlice["UDP"] > 0 {
		protocol = "all"
	} else if protocolsSlice["TCP"] > 0 {
		protocol = "tcp"
	} else if protocolsSlice["UDP"] > 0 {
		protocol = "udp"
	}
	// Gets persistence and Stickiness timeout in seconds
	persistence := ""
	if service.Spec.SessionAffinity == "ClientIP" {
		persistence = "srcip"
	} else if service.Spec.SessionAffinity == "None" {
		persistence = "none"
	}
	persistenceTTL := ""
	if service.Spec.SessionAffinityConfig != nil {
		if service.Spec.SessionAffinityConfig.ClientIP != nil {
			if service.Spec.SessionAffinityConfig.ClientIP.TimeoutSeconds != nil {
				// Value between 0 and 86400 seconds (1 day max)
				persistenceTTL = fmt.Sprint(*(service.Spec.SessionAffinityConfig.ClientIP.TimeoutSeconds))
			}
		}
	}

	// If there are no annotations, default values ​​are established
	mode := "snat"
	scheduler := ""
	helper := ""
	log := ""
	logprefix := ""
	// Get annotations
	var rgx = regexp.MustCompile(`[a-z]+$`)
	if service.ObjectMeta.Annotations != nil {
		for key, value := range service.ObjectMeta.Annotations {
			field := rgx.FindStringSubmatch(key)
			if strings.ToLower(string(field[0])) == "mode"{
				mode = value
			} else if strings.ToLower(string(field[0])) == "scheduler"{
				scheduler = value
			}else if strings.ToLower(string(field[0])) == "helper"{
				helper = value
			}else if strings.ToLower(string(field[0])) == "log"{
				log = value
			}else if strings.ToLower(string(field[0])) == "logprefix" && log != ""{
				logprefix = value
			}
		}
	}
	var farmsSlice []types.Farm
	serviceName := service.ObjectMeta.Name
	farmName := ""
	if len(service.Spec.Ports) > 1 {
		for _, port := range service.Spec.Ports {
			// Gets the name of the farm
			farmName = serviceName + "--" + port.Name
			// Gets the name of the protocol
			nameProtocol := strings.ToLower(string(port.Protocol))
			// Gets the port as string
			portString := fmt.Sprint(port.Port)
			// Fills the farm
			var farm = types.Farm{
				Name:         farmName,
				Family:       "ipv4",
				VirtualAddr:  service.Spec.ClusterIP,
				VirtualPorts: portString,
				Mode:         mode,
				Protocol:     nameProtocol,
				Scheduler:    scheduler,
				Helper:       helper,
				Log:          log,
				LogPrefix:    logprefix,
				State:        "up",
				Intraconnect: "on",
				Persistence:  fmt.Sprint(persistence),
				PersistTTL:   fmt.Sprint(persistenceTTL),
				Backends:     types.Backends{},
			}
			farmsSlice = append(farmsSlice, farm)
		}

	} else if len(service.Spec.Ports) == 1 {
		if service.Spec.Ports[0].Name == "" {
			farmName = serviceName + "--" + "default"
		} else {
			farmName = serviceName + "--" + service.Spec.Ports[0].Name
		}

		// Fills the farm
		var farm = types.Farm{
			Name:         farmName,
			Family:       "ipv4",
			VirtualAddr:  service.Spec.ClusterIP,
			VirtualPorts: ports,
			Mode:         mode,
			Protocol:     protocol,
			Scheduler:    scheduler,
			Helper:       helper,
			Log:          log,
			LogPrefix:    logprefix,
			State:        "up",
			Intraconnect: "on",
			Persistence:  fmt.Sprint(persistence),
			PersistTTL:   fmt.Sprint(persistenceTTL),
			Backends:     types.Backends{},
		}
		farmsSlice = append(farmsSlice, farm)
	}

	// Returns the filled struct
	return types.JSONnftlb{
		Farms: farmsSlice,
	}

}

// GetJSONnftlbFromEndpoints returns a JSONnftlb struct filled with any Endpoints data.
func GetJSONnftlbFromEndpoints(endpoints *v1.Endpoints) types.JSONnftlb {
	farmName := endpoints.ObjectMeta.Name
	portName := ""
	// Initializes farm/backends ID
	CreateFarmID(farmName)
	var farmsSlice []types.Farm
	for _, endpoint := range endpoints.Subsets {
		for _, address := range endpoint.Addresses {
			// Get ip of the backends
			ip := address.IP
			// Get name of the backends based on the deployment, if empty take name of the farm
			backendName := ""
			if address.TargetRef != nil {
				backendName = address.TargetRef.Name
			} else if address.TargetRef == nil {
				backendName = endpoints.ObjectMeta.Name
			}
			// Create backend for each port
			for _, port := range endpoint.Ports {
				var backends types.Backends
				var backend = types.Backend{
					Name:   fmt.Sprintf("%s", backendName),
					IPAddr: ip,
					State:  "up",
					Port:   fmt.Sprint(port.Port),
				}
				// Appends backend
				backends = append(backends, backend)

				if port.Name == "" {
					portName = "default"
				} else {
					portName = port.Name
				}

				var farm = types.Farm{
					Name:     fmt.Sprintf("%s%s%s", farmName, "--", portName),
					Backends: backends,
				}
				farmsSlice = append(farmsSlice, farm)
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
