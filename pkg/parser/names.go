package parser

import "fmt"

// FormatName returns a formatted name string for any nftlb object.
func FormatName(resourceName string, resourcePortName string) string {
	// The first possibility is the creation of one or several resources. If several are created from the same YAML
	// configuration file, we need to differentiate them (because they have the same resource name). For this case,
	// we make the resource name followed by the name of the resourcePort.
	// Example: "resource.Name + -- + resourcePort.Name" => "address--http"

	// The second possibility is when a single resource is created and some resourcePorts haven't been assigned a name.
	// It is assigned a default one called "default".
	// Example: "resource.Name + --default" => "address--default"

	if resourcePortName == "" {
		resourcePortName = "default"
	}

	return fmt.Sprintf("%s--%s", resourceName, resourcePortName)
}

// FormatNodePortName returns a formatted name (--nodePort suffix) for any nftlb object.
func FormatNodePortName(resourceName string, resourcePortName string) string {
	// The NodePort resource is called the same as the original resource by appending the string "nodePort".
	// Example: "address--http" => "address--http--nodeport".
	return fmt.Sprintf("%s--nodePort", FormatName(resourceName, resourcePortName))
}

// FormatExternalIPName returns a formatted name (--externalIP-index suffix) string for any nftlb object.
func FormatExternalIPName(resourceName string, resourcePortName string, index int) string {
	// The ExternalIP resource is called the same as the original resource by appending the string "externalIP-index".
	// Example: "address--http" => "address--http--externalIP-index".
	return fmt.Sprintf("%s--externalIP-%d", FormatName(resourceName, resourcePortName), index)
}
