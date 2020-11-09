package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	EndpointsChangesPending = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "kube_nftlb",
		Name:      "rules_endpoints_changes_pending",
		Help:      "How many Endpoints are pending to apply",
	})

	EndpointsChangesTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "kube_nftlb",
		Name:      "rules_endpoints_changes_total",
		Help:      "How many Endpoints changes have happened",
	})
)
