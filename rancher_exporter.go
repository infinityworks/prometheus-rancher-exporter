package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/infinityworksltd/prometheus-rancher-exporter/measure"
)

const (
	namespace = "rancher" // Used to prepand Prometheus metrics created by this exporter.
)

// Runtime variables, user controllable for targeting, authentication and filtering.
var (
	metricsPath   = getEnv("METRICS_PATH", "/metrics")       // Path under which to expose metrics
	listenAddress = getEnv("LISTEN_ADDRESS", ":9173")        // Address on which to expose metrics
	rancherURL    = os.Getenv("CATTLE_URL")                  // URL of Rancher Server API e.g. http://192.168.0.1:8080/v2-beta
	accessKey     = os.Getenv("CATTLE_ACCESS_KEY")           // Optional - Access Key for Rancher API
	secretKey     = os.Getenv("CATTLE_SECRET_KEY")           // Optional - Secret Key for Rancher API
	logLevel      = getEnv("LOG_LEVEL", "info")              // Optional - Set the logging level
	hideSys, _    = strconv.ParseBool(os.Getenv("HIDE_SYS")) // hideSys - Optional - Flag that indicates if the environment variable `HIDE_SYS` is set to a boolean true value
)

// Predefined variables that are used throughout the exporter
var (
	hostStates    = []string{"activating", "active", "deactivating", "error", "erroring", "inactive", "provisioned", "purged", "purging", "registering", "removed", "removing", "requested", "restoring", "updating_active", "updating_inactive"}
	stackStates   = []string{"activating", "active", "canceled_upgrade", "canceling_upgrade", "error", "erroring", "finishing_upgrade", "removed", "removing", "requested", "restarting", "rolling_back", "updating_active", "upgraded", "upgrading"}
	serviceStates = []string{"activating", "active", "canceled_upgrade", "canceling_upgrade", "deactivating", "finishing_upgrade", "inactive", "registering", "removed", "removing", "requested", "restarting", "rolling_back", "updating_active", "updating_inactive", "upgraded", "upgrading"}
	healthStates  = []string{"healthy", "unhealthy"}
	endpoints     = []string{"stacks", "services", "hosts"} // EndPoints the exporter will trawl
	stackRef      = make(map[string]string)                 // Stores the StackID and StackName as a map, used to provide label dimensions to service metrics

)

// getEnv - Allows us to supply a fallback option if nothing specified
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func main() {
	flag.Parse()

	// Sets the logging value for the exporter, defaults to info
	setLogLevel(logLevel)

	// check the rancherURL ($CATTLE_URL) has been provided correctly
	if rancherURL == "" {
		log.Fatal("CATTLE_URL must be set and non-empty")
	}

	log.Info("Starting Prometheus Exporter for Rancher")
	log.Info("Runtime Configuration in-use: URL of Rancher Server: ", rancherURL, " AccessKey: ", accessKey, "System Services Reported on: ", hideSys)

	// Register internal metrics used for tracking the exporter performance
	measure.Init()

	// Register a new Exporter
	Exporter := newExporter(rancherURL, accessKey, secretKey, hideSys)

	// Register Metrics from each of the endpoints
	// This invokes the Collect method through the prometheus client libraries.
	prometheus.MustRegister(Exporter)

	// Setup HTTP handler
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
	log.Printf("Starting Server on port %s and path %s", listenAddress, metricsPath)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
