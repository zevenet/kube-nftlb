package parser

var (
	// Active services and their variants
	farmsMap = map[string][]string{}

	// Active endpoints per service
	portsMap = map[string][]string{}

	// Active nodeports per service
	nodePortMap = map[string][]string{}

	// Active externalIPs per service
	externalIPsMap = map[string][]string{}

	// Maxconns per backend in a farm
	maxConnsMap = map[string]map[string]string{}
)
