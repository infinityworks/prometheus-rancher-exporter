package main

import (
	"fmt"
	"net/http"

	"github.com/fatih/structs"
	"github.com/infinityworks/prometheus-rancher-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	im "github.com/infinityworks/go-common/metrics"
	rm "github.com/infinityworks/prometheus-rancher-exporter/metrics"
)

var (
	applicationCfg config.Config
	rancherMetrics map[string]*prometheus.Desc
)

// Exporter Sets up all the runtime and metrics
type Exporter struct {
	rancherMetrics map[string]*prometheus.Desc
	config.Config
}

func init() {
	applicationCfg = config.Init()
	rancherMetrics = rm.Return()
}

func main() {
	log.WithFields(structs.Map(applicationCfg)).Info("Starting Prometheus Rancher Exporter")

	// Register internal metrics used for tracking the exporter performance
	im.Init()

	exporter := Exporter{
		rancherMetrics: rancherMetrics,
		Config:         applicationCfg,
	}

	// This invokes the Collect method through the prometheus client libraries.
	prometheus.MustRegister(&exporter)

	// Setup HTTP handler
	port := fmt.Sprintf(":%s", applicationCfg.ListenPort())
	http.Handle(applicationCfg.MetricsPath(), prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		                <head><title>Rancher Exporter</title></head>
		                <body>
		                   <h1>Rancher Exporter</h1>
		                   <p><a href='` + applicationCfg.MetricsPath() + `'>Metrics</a></p>
		                   </body>
		                </html>
		              `))
	})
	log.Fatal(http.ListenAndServe(port, nil))
}
