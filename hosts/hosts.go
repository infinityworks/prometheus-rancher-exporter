package hosts

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

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

// HostsData is used to store data from the hosts endpoint in the API
type HostsData struct {
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
func getJSONhosts(rancherURL string, accessKey string, secretKey string) (error, HostsData) {
	pulledData := HostsData{}
	client := &http.Client{}
	req, err := http.NewRequest("GET", rancherURL, nil)
	req.SetBasicAuth(accessKey, secretKey)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error Collecting JSON from API: ", err)
		panic(err)
	}
	return json.NewDecoder(resp.Body).Decode(&pulledData), pulledData

}

func (e *Exporter) scrapeHosts(rancherURL string, accessKey string, secretKey string, ch chan<- prometheus.Metric) error {

	for _, m := range e.gaugeVecs {
		m.Reset()
	}

	fmt.Println("Scraping: ", rancherURL+"/hosts/")
	err, hostsData := getJSONhosts(rancherURL+"/hosts/", accessKey, secretKey)
	if err != nil {
		return err
	}
	fmt.Println("JSON Fetched for hosts: ", hostsData)

	// Host Metrics
	for _, x := range hostsData.Data {

		// Set all the metrics to 0, unless we get a match

		e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "activating"}).Set(0)
		e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "active"}).Set(0)
		e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "deactivating"}).Set(0)
		e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "error"}).Set(0)
		e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "erroring"}).Set(0)
		e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "inactive"}).Set(0)
		e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "provisioned"}).Set(0)
		e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "purged"}).Set(0)
		e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "purging"}).Set(0)
		e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "registering"}).Set(0)
		e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "removed"}).Set(0)
		e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "removing"}).Set(0)
		e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "requested"}).Set(0)
		e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "restoring"}).Set(0)
		e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "updating_active"}).Set(0)
		e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "updating_inactive"}).Set(0)

		// Match states of the API to known values and override our values above.
		if x.State == "activating" {
			e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "activating"}).Set(1)
		} else if x.State == "active" {
			e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "active"}).Set(1)
		} else if x.State == "deactivating" {
			e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "deactivating"}).Set(1)
		} else if x.State == "error" {
			e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "error"}).Set(1)
		} else if x.State == "erroring" {
			e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "erroring"}).Set(1)
		} else if x.State == "inactive" {
			e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "inactive"}).Set(1)
		} else if x.State == "provisioned" {
			e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "provisioned"}).Set(1)
		} else if x.State == "purged" {
			e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "purged"}).Set(1)
		} else if x.State == "purging" {
			e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "purging"}).Set(1)
		} else if x.State == "registering" {
			e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "registering"}).Set(1)
		} else if x.State == "removed" {
			e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "removed"}).Set(1)
		} else if x.State == "removing" {
			e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "removing"}).Set(1)
		} else if x.State == "requested" {
			e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "requested"}).Set(1)
		} else if x.State == "restoring" {
			e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "restoring"}).Set(1)
		} else if x.State == "updating-active" {
			e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "updating_active"}).Set(1)
		} else if x.State == "updating-inactive" {
			e.gaugeVecs["HostState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname, "state": "updating_inactive"}).Set(1)
		}

	}

	return nil

}

// Collect function, called on by Prometheus Client
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	if err := e.scrapeHosts(e.rancherURL, e.accessKey, e.secretKey, ch); err != nil {
		log.Printf("Error scraping rancher url: %s", err)
		return
	}

	for _, m := range e.gaugeVecs {
		m.Collect(ch)
	}
}
