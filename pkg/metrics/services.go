package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ServicesChangesPending = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "kube_nftlb",
		Name:      "rules_services_changes_pending",
		Help:      "How many Services are pending to apply",
	})

	ServicesChangesTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "kube_nftlb",
		Name:      "rules_services_changes_total",
		Help:      "How many Services changes have happened",
	})
)
