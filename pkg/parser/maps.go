package parser

var (
	// Map [Service-Endpoints (name)] to []{ farms }
	farmsPerService = make(map[string][]string)

	// Map [farm (name)] to []{ backends }
	backendsPerFarm = make(map[string][]string)

	// Map [farm (name)] to []{ addresses }
	addressesPerFarm = make(map[string][]string)
)
