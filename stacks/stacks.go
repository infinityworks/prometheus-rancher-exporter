package stacks

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/infinityworksltd/prometheus-rancher-exporter/measure"
	"github.com/prometheus/client_golang/prometheus"
)

// Exporter collects Rancher stats from machine of a specified user and exports them using
// the prometheus metrics package.
type Exporter struct {
	rancherURL string
	accessKey  string
	secretKey  string
	mutex      sync.RWMutex
	gaugeVecs  map[string]*prometheus.GaugeVec
}

// Data is used to store data from the stacks endpoint in the API
type Data struct {
	Data []struct {
		HealthState string `json:"healthState"`
		Name        string `json:"name"`
		State       string `json:"state"`
	} `json:"data"`
}

//NewExporter creates the metrics we wish to monitor
func NewExporter(rancherURL string, accessKey string, secretKey string) *Exporter {

	gaugeVecs := make(map[string]*prometheus.GaugeVec)

	gaugeVecs["StackHealth"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "stack_health_status",
			Help:      "HealthState of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["StackState"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "stack_state",
			Help:      "State of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name", "state"})

	return &Exporter{
		gaugeVecs:  gaugeVecs,
		rancherURL: rancherURL,
		accessKey:  accessKey,
		secretKey:  secretKey,
	}
}

// Describe describes all the metrics ever exported by the Rancher exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {

	for _, m := range e.gaugeVecs {
		m.Describe(ch)
	}
}

// Gets the JSON response from the API and places it in the Data
func getJSON(collectionURL string, accessKey string, secretKey string) (error, Data) {

	start := time.Now()

	// Counter for internal exporter metrics
	measure.FunctionCountTotal.With(prometheus.Labels{"pkg": "stacks", "fnc": "getJSON"}).Inc()

	pulledData := Data{}
	client := &http.Client{}
	req, err := http.NewRequest("GET", collectionURL, nil)
	req.SetBasicAuth(accessKey, secretKey)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error Collecting JSON from API: ", err)
		panic(err)
	}

	// Timings recorded as part of internal metrics
	elapsed := float64((time.Since(start)) / time.Microsecond)
	measure.FunctionDurations.WithLabelValues("stacks", "getJSON").Observe(elapsed)

	return json.NewDecoder(resp.Body).Decode(&pulledData), pulledData

}

func (e *Exporter) gatherMetrics(rancherURL string, accessKey string, secretKey string, ch chan<- prometheus.Metric) error {

	for _, m := range e.gaugeVecs {
		m.Reset()
	}

	var stacksEndpoint string

	if strings.Contains(rancherURL, "v1") {
		fmt.Println("Version 1 API detected, using legacy API fields")
		stacksEndpoint = "/environments/"

	} else if strings.Contains(rancherURL, "v2") {
		fmt.Println("Version 2 API detected, using updated API fields")
		stacksEndpoint = "/stacks/"
	} else {
		fmt.Println("No known API version detected")
		stacksEndpoint = "/stacks/"
	}

	fmt.Println("Scraping: ", rancherURL+stacksEndpoint)
	err, Data := getJSON(rancherURL+stacksEndpoint, accessKey, secretKey)
	if err != nil {
		return err
	}
	fmt.Println("JSON Fetched for stacks: ", Data)

	// Stack Metrics
	for _, x := range Data.Data {

		var StackHealthState float64
		if x.HealthState == "healthy" {
			StackHealthState = 1
		}

		e.gaugeVecs["StackHealth"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackHealthState)

		// Pre-defines the known states from the Rancher API
		states := []string{"activating", "active", "canceled_upgrade", "canceling_upgrade", "error", "erroring", "finishing_upgrade", "removed", "removing", "requested", "restarting", "rolling_back", "updating_active", "upgraded", "upgrading"}

		// Set the state of the service to 1 when it matches one of the known states
		for _, y := range states {
			if x.State == y {
				e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": y}).Set(1)
			} else {
				e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": y}).Set(0)
			}
		}

	}
	return nil
}

// Collect function, called on by Prometheus Client
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	if err := e.gatherMetrics(e.rancherURL, e.accessKey, e.secretKey, ch); err != nil {
		log.Printf("Error scraping rancher url: %s", err)
		return
	}
	for _, m := range e.gaugeVecs {
		m.Collect(ch)
	}

}
