package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"

	"github.com/infinityworksltd/prometheus-rancher-exporter/hosts"
	"github.com/infinityworksltd/prometheus-rancher-exporter/services"
	"github.com/infinityworksltd/prometheus-rancher-exporter/stacks"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/infinityworksltd/prometheus-rancher-exporter/measure"
	"github.com/prometheus/log"
)

const (
	namespace = "rancher" // For Prometheus metrics.
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

var (
	metricsPath   = getenv("METRICS_PATH", "/metrics") // Path under which to expose metrics.
	listenAddress = getenv("LISTEN_ADDRESS", ":9010")  // Address on which to expose metrics.
	rancherURL    = os.Getenv("CATTLE_URL")            // URL of Rancher Server API e.g. http://192.168.0.1:8080/v2-beta
	accessKey     = os.Getenv("CATTLE_ACCESS_KEY")     // Optional - Access Key for Rancher API
	secretKey     = os.Getenv("CATTLE_SECRET_KEY")     // Optional - Secret Key for Rancher API

	// hideSys - Optional - Flag that indicates if the environment variable `HIDE_SYS` is set to a boolean true value.
	hideSys, _ = strconv.ParseBool(os.Getenv("HIDE_SYS"))
)

func main() {
	flag.Parse()
	if rancherURL == "" {
		log.Fatal("CATTLE_URL must be set and non-empty")
	}

	log.Info("Starting Prometheus Exporter for Rancher. Listen Address: ", listenAddress, " metricsPath: ", metricsPath, " rancherURL: ", rancherURL, " AccessKey: ", accessKey)
	log.Info("System Services Reported on:", hideSys)

	// Register internal metrics
	measure.Init()

	// Pass URL & Credentials out to the Exporters
	servicesExporter := services.NewExporter(rancherURL, accessKey, secretKey, hideSys)
	stacksExporter := stacks.NewExporter(rancherURL, accessKey, secretKey, hideSys)
	hostsExporter := hosts.NewExporter(rancherURL, accessKey, secretKey)

	// Register Metrics from each of the endpoints
	prometheus.MustRegister(servicesExporter)
	prometheus.MustRegister(stacksExporter)
	prometheus.MustRegister(hostsExporter)

	http.Handle(metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		                <head><title>Rancher exporter</title></head>
		                <body>
		                   <h1>rancher exporter</h1>
		                   <p><a href='` + metricsPath + `'>Metrics</a></p>
		                   </body>
		                </html>
		              `))
	})

	log.Infof("Starting Server: %s", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
