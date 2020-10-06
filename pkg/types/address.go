package types

// Address defines a nftlb address object. Equivalent to a k8s ServicePort.
type Address struct {
	Name     string `json:"name"`
	Family   string `json:"family"`
	IPAddr   string `json:"ip-addr"`
	Ports    string `json:"ports"`
	Protocol string `json:"protocol"`
}
