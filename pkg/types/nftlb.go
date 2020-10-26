package types

// Nftlb defines a struct made for nftlb requests.
type Nftlb struct {
	Addresses []Address `json:"addresses,omitempty"`
	Farms     []Farm    `json:"farms,omitempty"`
	Policies  []Policy  `json:"policies,omitempty"`
}
