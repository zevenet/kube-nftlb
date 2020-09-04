package parser

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

func getPersistence(service *corev1.Service) (string, string) {
	// First we get the persistence of our service. By default, annotations have priority ahead of the sessionAffinity and sessionAffinityConfig field.
	// If there are no annotations, the information in the sessionAffinity and sessionAffinityConfig field is collected.
	persistence := ""
	persistenceTTL := "60"
	var rgx = regexp.MustCompile(`[a-z]+$`)
	if service.ObjectMeta.Annotations != nil {
		for key, value := range service.ObjectMeta.Annotations {
			field := rgx.FindStringSubmatch(key)
			if strings.ToLower(string(field[0])) == "persistence" {
				//  check for multiple combination of values. Ej srcip+srcport && srcip+dstport
				splitField := strings.Split(value, "-")
				if len(splitField) > 1 {
					if splitField[0] == "srcip" && splitField[1] == "srcport" {
						persistence = splitField[0] + " " + splitField[1]
					} else if splitField[0] == "srcip" && splitField[1] == "dstport" {
						persistence = splitField[0] + " " + splitField[1]
					}
				} else if value == "srcip" || value == "dstip" || value == "srcport" || value == "dstport" || value == "srcmac" || value == "dstmac" {
					persistence = value
				}
			}
		}
	}
	if persistence == "" {
		if service.Spec.SessionAffinity == "ClientIP" {
			persistence = "srcip"
		} else if service.Spec.SessionAffinity == "None" {
			persistence = "none"
		}
	}
	if service.Spec.SessionAffinityConfig != nil {
		if service.Spec.SessionAffinityConfig.ClientIP != nil {
			if service.Spec.SessionAffinityConfig.ClientIP.TimeoutSeconds != nil {
				// Value between 0 and 86400 seconds (1 day max)
				persistenceTTL = fmt.Sprint(*(service.Spec.SessionAffinityConfig.ClientIP.TimeoutSeconds))
			}
		}
	}
	return persistence, persistenceTTL
}

func getAnnotations(service *corev1.Service, farmName string) (string, string, string, string, string, string) {
	// First try reading the annotations for fields that can be configured in the nftlb service
	// If there are no annotations for all the fields, default values ​​are set.
	// You don't need to worry about sending empty variables as it is configured so if a variable is sent empty it is not included in the json that configures the nftlb service.
	mode := "snat"
	scheduler := "rr"
	schedParam := "none"
	helper := ""
	log := ""
	logprefix := ""
	// We use a regular expression to filter the string and get the field to configure in the annotations
	// Always respecting the format of the string | service.kubernetes.io/kube-nftlb-load-balancer-X | where X is the field to configure
	var rgx = regexp.MustCompile(`[a-z]+$`)
	if service.ObjectMeta.Annotations != nil {
		for key, value := range service.ObjectMeta.Annotations {
			field := rgx.FindStringSubmatch(key)
			if strings.ToLower(string(field[0])) == "mode" {
				if value == "snat" || value == "dnat" || value == "dsr" || value == "stlsdnat" || value == "local" {
					mode = value
				}
			} else if strings.ToLower(string(field[0])) == "scheduler" {
				rgx = regexp.MustCompile(`^[a-z]+`)
				field = rgx.FindStringSubmatch(value)
				if value == "rr" {
					scheduler = value
				} else if value == "symhash" {
					scheduler = value
				} else if strings.ToLower(string(field[0])) == "hash" {
					splitField := strings.Split(value, "-")
					// check for multiple combination of values. Ej srcip+srcport
					if len(splitField) > 2 {
						if splitField[1] == "srcip" && splitField[2] == "srcport" {
							schedParam = splitField[1] + " " + splitField[2]
						}
					} else {
						valueHash := splitField[1]
						if valueHash == "srcip" || valueHash == "dstip" || valueHash == "srcport" || valueHash == "dstport" || valueHash == "srcmac" || valueHash == "dstmac" {
							schedParam = valueHash
						}
						scheduler = "hash"
					}

				}
			} else if strings.ToLower(string(field[0])) == "helper" {
				if value == "amanda" || value == "ftp" || value == "h323" || value == "irc" || value == "netbios-ns" || value == "pptp" || value == "sane" || value == "sip" || value == "snmp" || value == "tftp" {
					helper = value
				}
			} else if strings.ToLower(string(field[0])) == "log" {
				if value == "none" || value == "forward" || value == "output" {
					log = value
					logprefix = farmName
				}
			}
		}
	}
	return mode, scheduler, schedParam, helper, log, logprefix
}

func findMaxConns(service *corev1.Service) {
	var farmSlice []string
	backendMaxConnsMap := "0"
	serviceName := service.ObjectMeta.Name
	var rgx = regexp.MustCompile(`[a-z]+$`)
	if service.ObjectMeta.Annotations != nil {
		farmName := ""
		for _, port := range service.Spec.Ports {
			if port.Name == "" {
				farmName = assignFarmNameService(serviceName, "default")
			} else {
				farmName = assignFarmNameService(serviceName, port.Name)
			}
			farmSlice = append(farmSlice, farmName)
		}
		for key, value := range service.ObjectMeta.Annotations {
			field := rgx.FindStringSubmatch(key)
			if strings.ToLower(string(field[0])) == "maxconns" {
				backendMaxConnsMap = value
			}
		}
	}

	maxConnsMap[serviceName] = make(map[string]string)
	for _, farmName := range farmSlice {
		maxConnsMap[serviceName][farmName] = backendMaxConnsMap
	}
}

func findFamily(service *corev1.Service) string {
	// Find out what type of version the service IP has, by default the value ​​is ipv4
	family := "ipv4"
	localhostIp := net.ParseIP(service.Spec.ClusterIP)
	if localhostIp.To4() != nil {
		family = "ipv4"
	} else if localhostIp.To16() != nil {
		family = "ipv6"
	}
	return family
}

func assignFarmNameService(serviceName string, portName string) string {
	// We assign the name of the farm. Two possibilities are contemplated.
	// The first possibility is the creation of one or several services. If several are created from the same yaml configuration file we need to differentiate them (because they have the same service name). For this we add the name of the service followed by the name of the port
	// farmName = service.ObjectMeta.Name + "--" + port.Name

	// The second possibility is when a single service is created and it has not been assigned a port name. It is assigned a default one called "default"
	// farmName = service.ObjectMeta.Name + "--" + "default"
	farmName := serviceName + "--" + portName
	return farmName
}

func assignFarmNameNodePort(serviceName string, nodeportName string) string {
	// The nodeport service is called the same as the original service by adding the string node-port
	// Ej my-service--http, the nodeport service is called my-service--http--nodeport
	farmName := serviceName + "--" + nodeportName
	return farmName
}

func assignFarmNameExternalIPs(serviceName string, externalIPsName string) string {
	// The nodeport service is called the same as the original service by adding the string node-port
	// Ej my-service--http, the nodeport service is called my-service--http--externalIPsName
	farmName := serviceName + "--" + externalIPsName
	return farmName
}
