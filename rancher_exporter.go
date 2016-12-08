package main

import (
	"flag"
	"net/http"
	"os"
	"rancher-go-exporter/hosts"
	"rancher-go-exporter/services"
	"rancher-go-exporter/stacks"

	"fmt"

	"github.com/prometheus/client_golang/prometheus"
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
	metricsPath   = getenv("METRICS_PATH", "/metrics") //Path under which to expose metrics.
	listenAddress = getenv("LISTEN_ADDRESS", ":9010")  // Address on which to expose metrics.
	rancherURL    = os.Getenv("CATTLE_URL")            // URL of Rancher Server API
	accessKey     = os.Getenv("CATTLE_ACCESS_KEY")     // Access Key for Rancher API
	secretKey     = os.Getenv("CATTLE_SECRET_KEY")     // Secret Key for Rancher API
)

func main() {

	flag.Parse()
	fmt.Println("Starting Prometheus Exporter for Rancher. Listen Address: ", listenAddress, "metricsPath: ", metricsPath, "rancherURL: ", rancherURL, "AccessKey: ", accessKey)
	servicesExporter := services.NewExporter(rancherURL, accessKey, secretKey)
	stacksExporter := stacks.NewExporter(rancherURL, accessKey, secretKey)
	hostsExporter := hosts.NewExporter(rancherURL, accessKey, secretKey)
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
