package services

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
	hideSys    bool
	mutex      sync.RWMutex
	gaugeVecs  map[string]*prometheus.GaugeVec
}

// Data is used to store data from the services endpoint in the API
type Data struct {
	Data []struct {
		HealthState string `json:"healthState"`
		Name        string `json:"name"`
		Scale       int    `json:"scale"`
		State       string `json:"state"`
		System      bool   `json:"system"`
	} `json:"data"`
}

//NewExporter creates the metrics we wish to monitor
func NewExporter(rancherURL string, accessKey string, secretKey string, hideSys bool) *Exporter {

	gaugeVecs := make(map[string]*prometheus.GaugeVec)

	gaugeVecs["ServiceScale"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_scale",
			Help:      "scale of defined service as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["ServiceHealth"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_health_status",
			Help:      "HealthState of the service, as reported by the Rancher API. Either (1) or (0)",
		}, []string{"rancherURL", "name"})
	gaugeVecs["ServiceState"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_state",
			Help:      "State of the service, as reported by the Rancher API",
		}, []string{"rancherURL", "name", "state"})

	return &Exporter{
		gaugeVecs:  gaugeVecs,
		rancherURL: rancherURL,
		accessKey:  accessKey,
		secretKey:  secretKey,
		hideSys:    hideSys,
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
	measure.FunctionCountTotal.With(prometheus.Labels{"pkg": "services", "fnc": "getJSON"}).Inc()

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
	measure.FunctionDurations.WithLabelValues("services", "getJSON").Observe(elapsed)

	return json.NewDecoder(resp.Body).Decode(&pulledData), pulledData
}

func (e *Exporter) gatherMetrics(rancherURL string, accessKey string, secretKey string, hideSys bool, ch chan<- prometheus.Metric) error {

	for _, m := range e.gaugeVecs {
		m.Reset()
	}

	fmt.Println("Scraping: ", rancherURL+"/services/")
	err, Data := getJSON(rancherURL+"/services/", accessKey, secretKey)
	if err != nil {
		return err
	}
	fmt.Println("JSON Fetched for services: ", Data)

	// Service Metrics
	for _, x := range Data.Data {

		// If system services have been ignored, the loop simply skips them
		if hideSys == true && x.System == true {
			continue
		} else {

			// Get the Healthy State for a service
			if x.HealthState == "healthy" {
				e.gaugeVecs["ServiceHealth"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
			} else {
				e.gaugeVecs["ServiceHealth"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
			}

			// Set the scale of the service
			e.gaugeVecs["ServiceScale"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(float64(x.Scale))

			// Pre-defines the known states from the Rancher API
			states := []string{"activating", "active", "canceled_upgrade", "canceling_upgrade", "deactivasting", "finishing_upgrade", "inactive", "registering", "removed", "removing", "requested", "restarting", "rolling_back", "updating_active", "updating_inactive", "upgraded", "upgrading"}

			// Set the state of the service to 1 when it matches one of the known states
			for _, y := range states {
				if x.State == y {
					e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": y}).Set(1)
				} else {
					e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": y}).Set(0)
				}
			}
		}
	}

	return nil
}

// Collect function, called on by Prometheus Client
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	if err := e.gatherMetrics(e.rancherURL, e.accessKey, e.secretKey, e.hideSys, ch); err != nil {
		log.Printf("Error scraping rancher url: %s", err)
		return
	}

	for _, m := range e.gaugeVecs {
		m.Collect(ch)
	}
}
