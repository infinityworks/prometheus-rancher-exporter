package main

import (
	"flag"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/infinityworks/prometheus-rancher-exporter/measure"
)

const (
	namespace           = "rancher" // Used to prepand Prometheus metrics created by this exporter.
	defaultLabelsFilter = "^io.prometheus"
)

// Runtime variables, user controllable for targeting, authentication and filtering.
var (
	log = logrus.New()

	metricsPath   = getEnv("METRICS_PATH", "/metrics")            // Path under which to expose metrics
	listenAddress = getEnv("LISTEN_ADDRESS", ":9173")             // Address on which to expose metrics
	rancherURL    = os.Getenv("CATTLE_URL")                       // URL of Rancher Server API e.g. http://192.168.0.1:8080/v2-beta
	accessKey     = os.Getenv("CATTLE_ACCESS_KEY")                // Optional - Access Key for Rancher API
	secretKey     = os.Getenv("CATTLE_SECRET_KEY")                // Optional - Secret Key for Rancher API
	labelsFilter  = os.Getenv("LABELS_FILTER")                    // Optional - Filter for Rancher label names
	logLevel      = getEnv("LOG_LEVEL", "info")                   // Optional - Set the logging level
	resourceLimit = getEnv("API_LIMIT", "100")                    // Optional - Rancher API resource limit (default: 100)
	hideSys, _    = strconv.ParseBool(getEnv("HIDE_SYS", "true")) // hideSys - Optional - Flag that indicates if the environment variable `HIDE_SYS` is set to a boolean true value
)

// Predefined variables that are used throughout the exporter
var (
	agentStates   = []string{"activating", "active", "reconnecting", "disconnected", "disconnecting", "finishing-reconnect", "reconnected"}
	hostStates    = []string{"activating", "active", "deactivating", "disconnected", "error", "erroring", "inactive", "provisioned", "purged", "purging", "reconnecting", "registering", "removed", "removing", "requested", "restoring", "updating_active", "updating_inactive"}
	stackStates   = []string{"activating", "active", "canceled_upgrade", "canceling_upgrade", "error", "erroring", "finishing_upgrade", "removed", "removing", "requested", "restarting", "rolling_back", "updating_active", "upgraded", "upgrading"}
	serviceStates = []string{"activating", "active", "canceled_upgrade", "canceling_upgrade", "deactivating", "finishing_upgrade", "inactive", "registering", "removed", "removing", "requested", "restarting", "rolling_back", "updating_active", "updating_inactive", "upgraded", "upgrading"}
	healthStates  = []string{"healthy", "unhealthy", "initializing", "degraded", "started-once"}
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

	if labelsFilter == "" {
		labelsFilter = defaultLabelsFilter
	}

	labelsFilterRegexp, err := regexp.Compile(labelsFilter)
	if err != nil {
		log.Fatal("LABELS_FILTER must be valid regular expression")
	}

	log.Info("Starting Prometheus Exporter for Rancher")
	log.Info(
		"Runtime Configuration in-use: URL of Rancher Server: ",
		rancherURL,
		" Access key: ",
		accessKey,
		" System services hidden: ",
		hideSys,
		" Labels filter: ",
		labelsFilter,
	)

	// Register internal metrics used for tracking the exporter performance
	measure.Init()

	// Register a new Exporter
	exporter := newExporter(rancherURL, accessKey, secretKey, labelsFilterRegexp, hideSys, resourceLimit)

	// Register Metrics from each of the endpoints
	// This invokes the Collect method through the prometheus client libraries.
	prometheus.MustRegister(exporter)

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
