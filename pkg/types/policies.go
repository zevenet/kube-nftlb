package types

// Element defines what's inside a policy "elements" attribute.
type Element struct {
	Data string `json:"data"`
}

// Policy defines what's a policy nftlb object with its attributes.
type Policy struct {
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Family    string    `json:"family"`
	LogPrefix string    `json:"log-prefix,omitempty"`
	Elements  []Element `json:"elements"`
}

// Policies defines the "policies" nftlb object.
type Policies struct {
	Policies []Policy `json:"policies"`
}
