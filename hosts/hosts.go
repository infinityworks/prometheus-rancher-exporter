package hosts

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

// Data is used to store data from the hosts endpoint in the API
type Data struct {
	Data []struct {
		Hostname string `json:"hostname"`
		State    string `json:"state"`
	} `json:"data"`
}

// NewExporter returns an initialized Exporter.
func NewExporter(rancherURL string, accessKey string, secretKey string) *Exporter {

	gaugeVecs := make(map[string]*prometheus.GaugeVec)
	gaugeVecs["HostState"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      ("host_state"),
			Help:      "State of defined host as reported by the Rancher API",
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

// Gets the JSON response from the API and places it in the struct
func getJSON(rancherURL string, accessKey string, secretKey string) (error, Data) {

	start := time.Now()

	// Counter for internal exporter metrics
	measure.FunctionCountTotal.With(prometheus.Labels{"pkg": "hosts", "fnc": "getJSON"}).Inc()

	pulledData := Data{}
	client := &http.Client{}
	req, err := http.NewRequest("GET", rancherURL, nil)
	req.SetBasicAuth(accessKey, secretKey)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error Collecting JSON from API: ", err)
		panic(err)
	}

	// Timings recorded as part of internal metrics
	elapsed := float64((time.Since(start)) / time.Microsecond)
	measure.FunctionDurations.WithLabelValues("hosts", "getJSON").Observe(elapsed)

	return json.NewDecoder(resp.Body).Decode(&pulledData), pulledData

}

func (e *Exporter) gatherMetrics(rancherURL string, accessKey string, secretKey string, ch chan<- prometheus.Metric) error {

	for _, m := range e.gaugeVecs {
		m.Reset()
	}

	fmt.Println("Scraping: ", rancherURL+"/hosts/")
	err, Data := getJSON(rancherURL+"/hosts/", accessKey, secretKey)
	if err != nil {
		return err
	}
	fmt.Println("JSON Fetched for hosts: ", Data)

	// Host Metrics
	for _, x := range Data.Data {

		// Pre-defines the known states from the Rancher API
		states := []string{"activating", "active", "deactivating", "error", "erroring", "inactive", "provisioned", "purged", "purging", "registering", "removed", "removing", "requested", "restoring", "updating_active", "updating_inactive"}

		// Set the state of the service to 1 when it matches one of the known states
		for _, y := range states {
			if x.State == y {
				e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": y}).Set(1)
			} else {
				e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": y}).Set(0)
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
