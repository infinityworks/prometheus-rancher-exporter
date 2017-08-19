package main

import (
	"fmt"
	"net/http"

	"github.com/fatih/structs"
	"github.com/infinityworks/prometheus-rancher-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/infinityworks/prometheus-rancher-exporter/measure"
)

var (
	applicationCfg config.Config
	mets           map[string]*prometheus.Desc
)

// Exporter Sets up all the runtime and metrics
type Exporter struct {
	gaugeVecs map[string]*prometheus.Desc
	config.Config
}

func init() {
	applicationCfg = config.Init()
	mets = addMetrics()
}

func main() {
	log.WithFields(structs.Map(applicationCfg)).Info("Starting Prometheus Rancher Exporter")

	// Register internal metrics used for tracking the exporter performance
	measure.Init()

	exporter := Exporter{
		gaugeVecs: mets,
		Config:    applicationCfg,
	}

	// This invokes the Collect method through the prometheus client libraries.
	prometheus.MustRegister(&exporter)

	// Setup HTTP handler
	port := fmt.Sprintf(":%s", applicationCfg.ListenPort())
	http.Handle(applicationCfg.MetricsPath(), prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		                <head><title>Rancher exporter</title></head>
		                <body>
		                   <h1>rancher exporter</h1>
		                   <p><a href='` + applicationCfg.MetricsPath() + `'>Metrics</a></p>
		                   </body>
		                </html>
		              `))
	})
	log.Printf("Starting Server on port %s and path %s", applicationCfg.ListenPort(), applicationCfg.MetricsPath())
	log.Fatal(http.ListenAndServe(port, nil))
}

/*

TODO

- Add back in health status for services/stacks, talk to Bill
- Debug logging, fix passing this around. Put back in the initialisation bits
- Package Layout - simplify
- Add host metrics
- Add Stack metrics
- rename guagevecs
- Add any additional metrics that might be useful
- Update docs
- Push and release
- Add docker version to host metrics
- Add back in state for hosts
- Add kind check for host/stack
- Improve way we version metadata
-
*/
