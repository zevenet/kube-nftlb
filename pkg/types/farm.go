package types

// Backend defines a nftlb backend object with its properties. Equivalent to a k8s Pod.
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

// Farm defines a nftlb farm object with its properties. Equivalent to a k8s Service.
type Farm struct {
	Name         string    `json:"name"`
	Mode         string    `json:"mode,omitempty"`
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
	EstConnlimit string    `json:"est-connlimit,omitempty"`
	Backends     []Backend `json:"backends,omitempty"`
	Addresses    []Address `json:"addresses,omitempty"`
}
