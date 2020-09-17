package types

// DSR stores data to create a DSR interface.
type DSR struct {
	VirtualAddr  string
	VirtualPorts string
	DockerUID    []string
}
