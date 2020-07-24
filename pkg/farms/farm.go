package farms

func AssignFarmNameService(serviceName string, portName string) string {
	// We assign the name of the farm. Two possibilities are contemplated.
	// The first possibility is the creation of one or several services. If several are created from the same yaml configuration file we need to differentiate them (because they have the same service name). For this we add the name of the service followed by the name of the port
	// farmName = service.ObjectMeta.Name + "--" + port.Name

	// The second possibility is when a single service is created and it has not been assigned a port name. It is assigned a default one called "default"
	// farmName = service.ObjectMeta.Name + "--" + "default"
	farmName := serviceName + "--" + portName
	return farmName
}

func AssignFarmNameNodePort(serviceName string, nodeportName string) string {
	// The nodeport service is called the same as the original service by adding the string node-port
	// Ej my-service--http, the nodeport service is called my-service--http--nodeport
	farmName := serviceName + "--" + nodeportName
	return farmName
}
