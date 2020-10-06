package types

// ServiceData stores some useful values from a Service. The "Family" property must be found analyzing that Service.
type ServiceData struct {
	Name      string
	Family    string
	Type      string
	ClusterIP string
}
