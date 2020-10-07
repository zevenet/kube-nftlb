package parser

var (
	// Map [Service name] to []{ addresses names }
	addressesPerService = make(map[string][]string)

	// Map [Endpoints name] to []{ backends names }
	backendsPerEndpoints = make(map[string][]string)
)
