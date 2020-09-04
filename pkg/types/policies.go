package types

// Element defines what's inside a policy "elements" attribute.
type Element struct {
	Data string
}

// Policy defines what's a policy nftlb object with its attributes.
type Policy struct {
	Name      string
	Type      string
	Traffic   string
	Family    string
	LogPrefix string
	Elements  []Element
}

// Policies defines the "policies" nftlb object.
type Policies struct {
	Policies []Policy
}
