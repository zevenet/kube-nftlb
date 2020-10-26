package types

// DSR stores data to add or delete DSR interfaces.
type DSR struct {
	DockerUIDs   []string
	AddressesIPs []string
}
