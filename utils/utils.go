package utils

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/infinityworksltd/prometheus-rancher-exporter/measure"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

// getJson return json from server
func GetJson(url string, accessKey string, secretKey string, target interface{}) error {

	start := time.Now()

	// Counter for internal exporter metrics
	measure.FunctionCountTotal.With(prometheus.Labels{"pkg": "utils", "fnc": "GetJson"}).Inc()

	log.Info("Scraping: ", url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	req.SetBasicAuth(accessKey, secretKey)
	resp, err := client.Do(req)

	if err != nil {
		log.Error("Error Collecting JSON from API: ", err)
		panic(err)
	}

	// Timings recorded as part of internal metrics
	elapsed := float64((time.Since(start)) / time.Microsecond)
	measure.FunctionDurations.WithLabelValues("hosts", "getJSON").Observe(elapsed)

	return json.NewDecoder(resp.Body).Decode(target)
}

// StacksURLCheck - Checks the API version for Rancher to determine the correct URL
func StacksURLCheck(rancherURL string) string {

	var stacksEndpoint string

	if strings.Contains(rancherURL, "v1") {
		log.Info("Version 1 API detected, using legacy API fields")
		stacksEndpoint = (rancherURL + "/environments/")

	} else if strings.Contains(rancherURL, "v2") {
		log.Info("Version 2 API detected, using updated API fields")
		stacksEndpoint = (rancherURL + "/stacks/")

	} else {
		log.Info("No known API version detected, defaulting to /stacks/")
		stacksEndpoint = (rancherURL + "/stacks/")
	}

	return stacksEndpoint

}
