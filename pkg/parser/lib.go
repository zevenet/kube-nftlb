package parser

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/zevenet/kube-nftlb/pkg/config"
	"github.com/zevenet/kube-nftlb/pkg/log"
	"github.com/zevenet/kube-nftlb/pkg/types"

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
	maxConnsFarm := "0"

	var rgx = regexp.MustCompile(`[a-z]+$`)
	if service.ObjectMeta.Annotations == nil {
		log.WriteLog(types.DetailedLog, fmt.Sprintf("findMaxConns: No annotations found for Service %s", service.Name))
	}

	for key, value := range service.ObjectMeta.Annotations {
		field := rgx.FindStringSubmatch(key)
		if strings.ToLower(string(field[0])) == "maxconns" {
			maxConnsFarm = value
		}
	}

	for _, servicePort := range service.Spec.Ports {
		farmName := FormatFarmName(service.Name, servicePort.Name)
		maxConnsMap[farmName] = maxConnsFarm
	}
}

func findFamily(service *corev1.Service) string {
	if localhostIP := net.ParseIP(service.Spec.ClusterIP); localhostIP.To4() != nil {
		return "ipv4"
	}
	return "ipv6"
}

func findIface(mode string) string {
	if mode == "dsr" {
		return config.DockerInterfaceBridge
	}
	return ""
}

func findProtocol(servicePort *corev1.ServicePort) string {
	return strings.ToLower(string(servicePort.Protocol))
}

func findVirtualPorts(servicePort *corev1.ServicePort) string {
	return strconv.FormatInt(int64(servicePort.Port), 10)
}

func findVirtualPortsNodePort(servicePort *corev1.ServicePort) string {
	return strconv.FormatInt(int64(servicePort.NodePort), 10)
}

// FormatFarmName returns a formatted farm name string for nftlb (regular Service).
func FormatFarmName(resourceName string, resourcePortName string) string {
	// The first possibility is the creation of one or several resources. If several are created from the same YAML
	// configuration file, we need to differentiate them (because they have the same resource name). For this we make
	// the resource name followed by the name of the resourcePort.
	// Example: "resource.Name + -- + resourcePort.Name" => "farm--http"

	// The second possibility is when a single resource is created and some resourcePorts haven't been assigned a name.
	// It is assigned a default one called "default".
	// Example: "resource.Name + --default" => "farm--default"

	if resourcePortName == "" {
		resourcePortName = "default"
	}

	return fmt.Sprintf("%s--%s", resourceName, resourcePortName)
}

// FormatNodePortFarmName returns a formatted farm name (--nodePort suffix) for nftlb.
func FormatNodePortFarmName(resourceName string, resourcePortName string) string {
	// The NodePort resource is called the same as the original Service by appending the string "nodePort".
	// Example: "farm--http" => "farm--http--nodeport".
	return fmt.Sprintf("%s--nodePort", FormatFarmName(resourceName, resourcePortName))
}

// FormatExternalIPFarmName returns a formatted farm name (--externalIP-index suffix) string for nftlb.
func FormatExternalIPFarmName(resourceName string, resourcePortName string, index int) string {
	// The ExternalIP resource is called the same as the original Service by appending the string "externalIP-index".
	// Example: "farm--http" => "farm--http--externalIP-index".
	return fmt.Sprintf("%s--externalIP-%d", FormatFarmName(resourceName, resourcePortName), index)
}
