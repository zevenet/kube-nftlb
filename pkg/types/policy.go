package types

// Element defines what's inside a policy "elements" attribute.
type Element struct {
	Data string `json:"data"`
}

// Policy defines a nftlb policy object with its properties.
type Policy struct {
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Family    string    `json:"family"`
	LogPrefix string    `json:"log-prefix,omitempty"`
	Elements  []Element `json:"elements"`
}
