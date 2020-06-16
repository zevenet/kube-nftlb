package types

// Backend defines any backend with its properties.
type Backend struct {
	Name     string `json:"name"`
	IPAddr   string `json:"ip-addr"`
	Weight   string `json:"weight,omitempty"`
	Priority string `json:"priority,omitempty"`
	Mark     string `json:"mark,omitempty"`
	State    string `json:"state,omitempty"`
	Port	 string `json:"port,omitempty"`
}

// Backends defines a group of backends in any farm.
type Backends []Backend

// Farm defines any farm with its properties.
type Farm struct {
	Name         string   `json:"name"`
	Family       string   `json:"family,omitempty"`
	VirtualAddr  string   `json:"virtual-addr,omitempty"`
	VirtualPorts string   `json:"virtual-ports,omitempty"`
	Mode         string   `json:"mode,omitempty"`
	Protocol     string   `json:"protocol,omitempty"`
	Scheduler    string   `json:"scheduler,omitempty"`
	Helper       string   `json:"helper,omitempty"`
	Log          string   `json:"log,omitempty"`
	Mark         string   `json:"mark,omitempty"`
	Priority     string   `json:"priority,omitempty"`
	State        string   `json:"state,omitempty"`
	Intraconnect string   `json:"intra-connect,omitempty"`
	Persistence  string   `json:"persistence,omitempty"`
        PersistTTL   string   `json:"persist-ttl,omitempty"`
	Backends     Backends `json:"backends"`

}

// Farms defines a group of farms.
type Farms []Farm

// JSONnftlb is a JSON object made for nftlb requests.
type JSONnftlb struct {
	Farms Farms `json:"farms"`
}
