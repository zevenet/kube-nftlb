package parser

var (
	// Map [(Node|Port) Service] to []{ farm names }
	farmsPerServiceMap = map[string][]string{}

	// Map [ExternalIP (Service|Endpoints)] to []{ farm names }
	farmsPerExternalIPResourceMap = map[string][]string{}

	// Map [Endpoints] to []{ farm names }
	farmsPerEndpointMap = map[string][]string{}

	// Map [Endpoints (farms)] to []{ backend names }
	backendsPerFarm = map[string][]string{}

	// Map [backend name] to "maxconn" value
	maxConnsMap = map[string]string{}
)

func existsFarm(serviceName string, expectedFarmName string) bool {
	for _, farmName := range farmsPerServiceMap[serviceName] {
		if farmName == expectedFarmName {
			return true
		}
	}
	return false
}
