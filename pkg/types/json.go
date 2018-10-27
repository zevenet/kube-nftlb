package types

// Backends defines a group of backends in any farm.
type Backends []struct {
	Name     string `json:"name"`
	IPAddr   string `json:"ip-addr"`
	Weight   string `json:"weight,omitempty"`
	Priority string `json:"priority,omitempty"`
	State    string `json:"state,omitempty"`
}

// Farms defines a group of farms.
type Farms []struct {
	Name         string   `json:"name"`
	Iface        string   `json:"iface,omitempty"`
	Oface        string   `json:"oface,omitempty"`
	Family       string   `json:"family"`
	EtherAddr    string   `json:"ether-addr,omitempty"`
	VirtualAddr  string   `json:"virtual-addr"`
	VirtualPorts string   `json:"virtual-ports"`
	Mode         string   `json:"mode"`
	Protocol     string   `json:"protocol,omitempty"`
	Scheduler    string   `json:"scheduler,omitempty"`
	Helper       string   `json:"helper,omitempty"`
	Log          string   `json:"log,omitempty"`
	Priority     string   `json:"priority,omitempty"`
	State        string   `json:"state,omitempty"`
	Backends     Backends `json:"backends"`
}

// JSONnftlb is a JSON object made for nftlb requests.
type JSONnftlb struct {
	Farms Farms `json:"farms"`
}
