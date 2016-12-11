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

	gaugeVecs["HostStateActivating"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      ("host_state_activating"),
			Help:      "State of defined host as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["HostStateActive"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      ("host_state_active"),
			Help:      "State of defined host as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["HostStateDeactivating"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      ("host_state_deactivating"),
			Help:      "State of defined host as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["HostStateError"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      ("host_state_error"),
			Help:      "State of defined host as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["HostStateErroring"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      ("host_state_erroring"),
			Help:      "State of defined host as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["HostStateInactive"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      ("host_state_inactive"),
			Help:      "State of defined host as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["HostStateProvisioned"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      ("host_state_provisioned"),
			Help:      "State of defined host as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["HostStatePurged"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      ("host_state_purged"),
			Help:      "State of defined host as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["HostStatePurging"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      ("host_state_purging"),
			Help:      "State of defined host as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["HostStateRegistering"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      ("host_state_registering"),
			Help:      "State of defined host as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["HostStateRemoved"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      ("host_state_removed"),
			Help:      "State of defined host as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["HostStateRemoving"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      ("host_state_removing"),
			Help:      "State of defined host as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["HostStateRequested"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      ("host_state_requested"),
			Help:      "State of defined host as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["HostStateRestoring"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      ("host_state_restoring"),
			Help:      "State of defined host as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["HostStateUpdatingActive"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      ("host_state_updating_active"),
			Help:      "State of defined host as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["HostStateUpdatingInactive"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      ("host_state_updating_inactive"),
			Help:      "State of defined host as reported by Rancher",
		}, []string{"rancherURL", "name"})

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

		e.gaugeVecs["HostStateActivating"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(0)
		e.gaugeVecs["HostStateActive"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(0)
		e.gaugeVecs["HostStateDeactivating"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(0)
		e.gaugeVecs["HostStateError"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(0)
		e.gaugeVecs["HostStateErroring"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(0)
		e.gaugeVecs["HostStateInactive"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(0)
		e.gaugeVecs["HostStateProvisioned"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(0)
		e.gaugeVecs["HostStatePurged"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(0)
		e.gaugeVecs["HostStatePurging"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(0)
		e.gaugeVecs["HostStateRegistering"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(0)
		e.gaugeVecs["HostStateRemoved"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(0)
		e.gaugeVecs["HostStateRemoving"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(0)
		e.gaugeVecs["HostStateRequested"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(0)
		e.gaugeVecs["HostStateRestoring"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(0)
		e.gaugeVecs["HostStateUpdatingActive"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(0)
		e.gaugeVecs["HostStateUpdatingInactive"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(0)

		// Match states of the API to known values and override our values above.
		if x.State == "activating" {
			e.gaugeVecs["HostStateActivating"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(1)
		} else if x.State == "active" {
			e.gaugeVecs["HostStateActive"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(1)
		} else if x.State == "deactivating" {
			e.gaugeVecs["HostStateDeactivating"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(1)
		} else if x.State == "error" {
			e.gaugeVecs["HostStateError"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(1)
		} else if x.State == "erroring" {
			e.gaugeVecs["HostStateErroring"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(1)
		} else if x.State == "inactive" {
			e.gaugeVecs["HostStateInactive"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(1)
		} else if x.State == "provisioned" {
			e.gaugeVecs["HostStateProvisioned"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(1)
		} else if x.State == "purged" {
			e.gaugeVecs["HostStatePurged"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(1)
		} else if x.State == "purging" {
			e.gaugeVecs["HostStatePurging"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(1)
		} else if x.State == "registering" {
			e.gaugeVecs["HostStateRegistering"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(1)
		} else if x.State == "removed" {
			e.gaugeVecs["HostStateRemoved"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(1)
		} else if x.State == "removing" {
			e.gaugeVecs["HostStateRemoving"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(1)
		} else if x.State == "requested" {
			e.gaugeVecs["HostStateRequested"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(1)
		} else if x.State == "restoring" {
			e.gaugeVecs["HostStateRestoring"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(1)
		} else if x.State == "updating-active" {
			e.gaugeVecs["HostStateUpdatingActive"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(1)
		} else if x.State == "updating-inactive" {
			e.gaugeVecs["HostStateUpdatingInactive"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(1)
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
