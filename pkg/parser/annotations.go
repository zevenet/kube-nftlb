package parser

import (
	"regexp"
	"strconv"

	"github.com/zevenet/kube-nftlb/pkg/types"

	corev1 "k8s.io/api/core/v1"
)

// A regular expression is used to filter the annotation and get the field to configure,
// always respecting the format of the annotation "service.kubernetes.io/kube-nftlb-load-balancer-X".
var rgxAnnotations = regexp.MustCompile(`^service.kubernetes.io/kube-nftlb-load-balancer-`)

func getAnnotations(service *corev1.Service) *types.Annotations {
	// Default values
	annotations := &types.Annotations{
		Persistence:  "none",
		PersistTTL:   "60",
		Mode:         "snat",
		Scheduler:    "rr",
		SchedParam:   "none",
		Helper:       "none",
		Log:          "none",
		EstConnlimit: "0",
	}

	// Override default Persistence value if SessionAffinity is defined as "ClientIP"
	if service.Spec.SessionAffinity == "ClientIP" {
		annotations.Persistence = "srcip"
	}

	// Override default PersistTTL value if SessionAffinityConfig has TimeoutSeconds
	if affConfig := service.Spec.SessionAffinityConfig; affConfig != nil && affConfig.ClientIP != nil && affConfig.ClientIP.TimeoutSeconds != nil {
		// Value between 0 and 86400 seconds (1 day at most)
		annotations.PersistTTL = strconv.FormatInt(int64(*affConfig.ClientIP.TimeoutSeconds), 32)
	}

	// Read every annotation from this Service
	for key, value := range service.ObjectMeta.Annotations {
		// Match annotation key against the regex and remove the matched regex text
		match := rgxAnnotations.ReplaceAllString(key, "")

		// TODO Update README with changes made to annotations (more akin to how nftlb accepts values)
		switch match {
		case "mode":
			annotations.Mode = value
		case "persistence":
			annotations.Persistence = value
		case "scheduler":
			annotations.Scheduler = value
		case "sched-param":
			annotations.SchedParam = value
		case "helper":
			annotations.Helper = value
		case "est-connlimit":
			annotations.EstConnlimit = value
		case "log":
			annotations.Log = value
			annotations.LogPrefix = service.Name
		}
	}

	annotations.Iface = findIface(annotations.Mode)

	return annotations
}
