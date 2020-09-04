package types

// Backend defines any backend with its properties.
type Backend struct {
	Name         string `json:"name"`
	IPAddr       string `json:"ip-addr"`
	Weight       string `json:"weight,omitempty"`
	Priority     string `json:"priority,omitempty"`
	Mark         string `json:"mark,omitempty"`
	State        string `json:"state,omitempty"`
	Port         string `json:"port,omitempty"`
	EstConnlimit string `json:"est-connlimit,omitempty"`
}

// Farm defines any farm with its properties.
type Farm struct {
	Name         string    `json:"name"`
	Family       string    `json:"family,omitempty"`
	VirtualAddr  string    `json:"virtual-addr,omitempty"`
	VirtualPorts string    `json:"virtual-ports,omitempty"`
	Mode         string    `json:"mode,omitempty"`
	Protocol     string    `json:"protocol,omitempty"`
	Scheduler    string    `json:"scheduler,omitempty"`
	SchedParam   string    `json:"sched-param,omitempty"`
	Helper       string    `json:"helper,omitempty"`
	Log          string    `json:"log,omitempty"`
	LogPrefix    string    `json:"log-prefix,omitempty"`
	Mark         string    `json:"mark,omitempty"`
	Priority     string    `json:"priority,omitempty"`
	State        string    `json:"state,omitempty"`
	IntraConnect string    `json:"intra-connect,omitempty"`
	Persistence  string    `json:"persistence,omitempty"`
	PersistTTL   string    `json:"persist-ttl,omitempty"`
	Iface        string    `json:"iface,omitempty"`
	Backends     []Backend `json:"backends"`
}

// Farms defines a struct made for nftlb requests.
type Farms struct {
	Farms []Farm `json:"farms"`
}
