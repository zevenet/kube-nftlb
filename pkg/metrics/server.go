package metrics

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	defaultAddress = ":9195"
	address        = flag.String("listen-address", defaultAddress, "The address to listen on for HTTP requests.")
	collectors     = []prometheus.Collector{
		EndpointsChangesPending,
		EndpointsChangesTotal,
		ServicesChangesPending,
		ServicesChangesTotal,
	}
)

func StartServer() {
	flag.Parse()

	if *address == "" {
		*address = defaultAddress
	}

	for _, collector := range collectors {
		prometheus.MustRegister(collector)
	}

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>kube-nftlb Exporter</title></head>
             <body>
             <h1>kube-nftlb Exporter</h1>
             <p><a href='/metrics'>Metrics</a></p>
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(*address, nil))
}
