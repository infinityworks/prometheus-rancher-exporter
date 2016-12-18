package stacks

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

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

// StacksData is used to store data from the stacks endpoint in the API
type StacksData struct {
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

// Gets the JSON response from the API and places it in the stacksData
func getJSONstacks(collectionURL string, accessKey string, secretKey string) (error, StacksData) {

	// Counter for internal exporter metrics
	measure.FunctionCountTotal.With(prometheus.Labels{"pkg": "stacks", "fnc": "getJSONstacks"}).Inc()

	pulledData := StacksData{}
	client := &http.Client{}
	req, err := http.NewRequest("GET", collectionURL, nil)
	req.SetBasicAuth(accessKey, secretKey)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error Collecting JSON from API: ", err)
		panic(err)
	}
	return json.NewDecoder(resp.Body).Decode(&pulledData), pulledData

}

func (e *Exporter) scrapeStacks(rancherURL string, accessKey string, secretKey string, ch chan<- prometheus.Metric) error {

	// Counter for internal exporter metrics
	measure.FunctionCountTotal.With(prometheus.Labels{"pkg": "stacks", "fnc": "scrapeStacks"}).Inc()

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
	err, stacksData := getJSONstacks(rancherURL+stacksEndpoint, accessKey, secretKey)
	if err != nil {
		return err
	}
	fmt.Println("JSON Fetched for stacks: ", stacksData)

	// Stack Metrics
	for _, x := range stacksData.Data {

		var StackHealthState float64
		if x.HealthState == "healthy" {
			StackHealthState = 1
		}
		e.gaugeVecs["StackHealth"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackHealthState)

		// Set all the metrics to 0, unless we get a match
		e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "activating"}).Set(0)
		e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "active"}).Set(0)
		e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "canceled_upgrade"}).Set(0)
		e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "canceling_upgrade"}).Set(0)
		e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "error"}).Set(0)
		e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "erroring"}).Set(0)
		e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "finishing_upgrade"}).Set(0)
		e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "removed"}).Set(0)
		e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "removing"}).Set(0)
		e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "requested"}).Set(0)
		e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "rolling_back"}).Set(0)
		e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "updating_active"}).Set(0)
		e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "upgraded"}).Set(0)
		e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "upgrading"}).Set(0)

		// Match states of the API to known values and override our values above.
		if x.State == "activating" {
			e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "activating"}).Set(1)
		} else if x.State == "active" {
			e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "active"}).Set(1)
		} else if x.State == "canceled-upgrade" {
			e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "canceled_upgrade"}).Set(1)
		} else if x.State == "canceling-upgrade" {
			e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "canceling_upgrade"}).Set(1)
		} else if x.State == "error" {
			e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "error"}).Set(1)
		} else if x.State == "erroring" {
			e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state:": "erroring"}).Set(1)
		} else if x.State == "finishing-upgrade" {
			e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "finishing_upgrade"}).Set(1)
		} else if x.State == "removed" {
			e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "removed"}).Set(1)
		} else if x.State == "removing" {
			e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "removing"}).Set(1)
		} else if x.State == "requested" {
			e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "requested"}).Set(1)
		} else if x.State == "rolling-back" {
			e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "rolling_back"}).Set(1)
		} else if x.State == "updating-active" {
			e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "updating_active"}).Set(1)
		} else if x.State == "upgraded" {
			e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "upgraded"}).Set(1)
		} else if x.State == "upgrading" {
			e.gaugeVecs["StackState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "upgrading"}).Set(1)

		}

	}
	return nil
}

// Collect function, called on by Prometheus Client
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	// Counter for internal exporter metrics
	measure.FunctionCountTotal.With(prometheus.Labels{"pkg": "stacks", "fnc": "Collect"}).Inc()

	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	if err := e.scrapeStacks(e.rancherURL, e.accessKey, e.secretKey, ch); err != nil {
		log.Printf("Error scraping rancher url: %s", err)
		return
	}
	for _, m := range e.gaugeVecs {
		m.Collect(ch)
	}
}
