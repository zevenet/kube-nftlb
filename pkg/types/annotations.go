package types

// Annotations stores values that can be passed to nftlb through k8s annotations.
type Annotations struct {
	Persistence  string
	PersistTTL   string
	Mode         string
	Scheduler    string
	SchedParam   string
	Helper       string
	Log          string
	LogPrefix    string
	EstConnlimit string
}
