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
	rancherURL                string
	accessKey                 string
	secretKey                 string
	mutex                     sync.RWMutex
	HostStateActivating       *prometheus.GaugeVec
	HostStateActive           *prometheus.GaugeVec
	HostStateDeactivating     *prometheus.GaugeVec
	HostStateError            *prometheus.GaugeVec
	HostStateErroring         *prometheus.GaugeVec
	HostStateInactive         *prometheus.GaugeVec
	HostStateProvisioned      *prometheus.GaugeVec
	HostStatePurged           *prometheus.GaugeVec
	HostStatePurging          *prometheus.GaugeVec
	HostStateRegistering      *prometheus.GaugeVec
	HostStateRemoved          *prometheus.GaugeVec
	HostStateRemoving         *prometheus.GaugeVec
	HostStateRequested        *prometheus.GaugeVec
	HostStateRestoring        *prometheus.GaugeVec
	HostStateUpdatingActive   *prometheus.GaugeVec
	HostStateUpdatingInactive *prometheus.GaugeVec
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
	return &Exporter{
		rancherURL: rancherURL,
		accessKey:  accessKey,
		secretKey:  secretKey,
		HostStateActivating: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      ("host_state_activating"),
				Help:      "State of defined host as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		HostStateActive: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      ("host_state_active"),
				Help:      "State of defined host as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		HostStateDeactivating: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      ("host_state_deactivating"),
				Help:      "State of defined host as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		HostStateError: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      ("host_state_error"),
				Help:      "State of defined host as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		HostStateErroring: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      ("host_state_erroring"),
				Help:      "State of defined host as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		HostStateInactive: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      ("host_state_inactive"),
				Help:      "State of defined host as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		HostStateProvisioned: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      ("host_state_provisioned"),
				Help:      "State of defined host as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		HostStatePurged: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      ("host_state_purged"),
				Help:      "State of defined host as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		HostStatePurging: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      ("host_state_purging"),
				Help:      "State of defined host as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		HostStateRegistering: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      ("host_state_registering"),
				Help:      "State of defined host as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		HostStateRemoved: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      ("host_state_removed"),
				Help:      "State of defined host as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		HostStateRemoving: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      ("host_state_removing"),
				Help:      "State of defined host as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		HostStateRequested: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      ("host_state_requested"),
				Help:      "State of defined host as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		HostStateRestoring: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      ("host_state_restoring"),
				Help:      "State of defined host as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		HostStateUpdatingActive: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      ("host_state_updating_active"),
				Help:      "State of defined host as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		HostStateUpdatingInactive: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      ("host_state_updating_inactive"),
				Help:      "State of defined host as reported by Rancher",
			}, []string{"rancherURL", "name"}),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.HostStateActivating.Describe(ch)
	e.HostStateActive.Describe(ch)
	e.HostStateDeactivating.Describe(ch)
	e.HostStateError.Describe(ch)
	e.HostStateErroring.Describe(ch)
	e.HostStateInactive.Describe(ch)
	e.HostStateProvisioned.Describe(ch)
	e.HostStatePurged.Describe(ch)
	e.HostStatePurging.Describe(ch)
	e.HostStateRegistering.Describe(ch)
	e.HostStateRemoved.Describe(ch)
	e.HostStateRemoving.Describe(ch)
	e.HostStateRequested.Describe(ch)
	e.HostStateRestoring.Describe(ch)
	e.HostStateUpdatingActive.Describe(ch)
	e.HostStateUpdatingInactive.Describe(ch)
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
	e.HostStateActivating.Reset()
	e.HostStateActive.Reset()
	e.HostStateDeactivating.Reset()
	e.HostStateError.Reset()
	e.HostStateErroring.Reset()
	e.HostStateInactive.Reset()
	e.HostStateProvisioned.Reset()
	e.HostStatePurged.Reset()
	e.HostStatePurging.Reset()
	e.HostStateRegistering.Reset()
	e.HostStateRemoved.Reset()
	e.HostStateRemoving.Reset()
	e.HostStateRequested.Reset()
	e.HostStateRestoring.Reset()
	e.HostStateUpdatingActive.Reset()
	e.HostStateUpdatingInactive.Reset()

	fmt.Println("Scraping: ", rancherURL+"hosts/")
	err, hostsData := getJSONhosts(rancherURL+"hosts/", accessKey, secretKey)
	if err != nil {
		return err
	}
	fmt.Println("JSON Fetched for hosts: ", hostsData)

	// Host Metrics
	for _, x := range hostsData.Data {

		var HostStateActivating float64
		var HostStateActive float64
		var HostStateDeactivating float64
		var HostStateError float64
		var HostStateErroring float64
		var HostStateInactive float64
		var HostStateProvisioned float64
		var HostStatePurged float64
		var HostStatePurging float64
		var HostStateRegistering float64
		var HostStateRemoved float64
		var HostStateRemoving float64
		var HostStateRequested float64
		var HostStateRestoring float64
		var HostStateUpdatingActive float64
		var HostStateUpdatingInactive float64

		if x.State == "activating" {
			HostStateActivating = 1
		} else if x.State == "active" {
			HostStateActive = 1
		} else if x.State == "deactivating" {
			HostStateDeactivating = 1
		} else if x.State == "error" {
			HostStateError = 1
		} else if x.State == "erroring" {
			HostStateErroring = 1
		} else if x.State == "inactive" {
			HostStateInactive = 1
		} else if x.State == "provisioned" {
			HostStateProvisioned = 1
		} else if x.State == "purged" {
			HostStatePurged = 1
		} else if x.State == "purging" {
			HostStatePurging = 1
		} else if x.State == "registering" {
			HostStateRegistering = 1
		} else if x.State == "removed" {
			HostStateRemoved = 1
		} else if x.State == "removing" {
			HostStateRemoving = 1
		} else if x.State == "requested" {
			HostStateRequested = 1
		} else if x.State == "restoring" {
			HostStateRestoring = 1
		} else if x.State == "updating-active" {
			HostStateUpdatingActive = 1
		} else if x.State == "updating-inactive" {
			HostStateUpdatingInactive = 1
		}

		e.HostStateActivating.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(HostStateActivating)
		e.HostStateActive.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(HostStateActive)
		e.HostStateDeactivating.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(HostStateDeactivating)
		e.HostStateError.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(HostStateError)
		e.HostStateErroring.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(HostStateErroring)
		e.HostStateInactive.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(HostStateInactive)
		e.HostStateProvisioned.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(HostStateProvisioned)
		e.HostStatePurged.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(HostStatePurged)
		e.HostStatePurging.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(HostStatePurging)
		e.HostStateRegistering.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(HostStateRegistering)
		e.HostStateRemoved.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(HostStateRemoved)
		e.HostStateRemoving.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(HostStateRemoving)
		e.HostStateRequested.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(HostStateRequested)
		e.HostStateRestoring.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(HostStateRestoring)
		e.HostStateUpdatingActive.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(HostStateUpdatingActive)
		e.HostStateUpdatingInactive.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Hostname}).Set(HostStateUpdatingInactive)

	}

	return nil

}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	if err := e.scrapeHosts(e.rancherURL, e.accessKey, e.secretKey, ch); err != nil {
		log.Printf("Error scraping rancher url: %s", err)
		return
	}

	e.HostStateActivating.Collect(ch)
	e.HostStateActive.Collect(ch)
	e.HostStateDeactivating.Collect(ch)
	e.HostStateError.Collect(ch)
	e.HostStateErroring.Collect(ch)
	e.HostStateInactive.Collect(ch)
	e.HostStateProvisioned.Collect(ch)
	e.HostStatePurged.Collect(ch)
	e.HostStatePurging.Collect(ch)
	e.HostStateRegistering.Collect(ch)
	e.HostStateRemoved.Collect(ch)
	e.HostStateRemoving.Collect(ch)
	e.HostStateRequested.Collect(ch)
	e.HostStateRestoring.Collect(ch)
	e.HostStateUpdatingActive.Collect(ch)
	e.HostStateUpdatingInactive.Collect(ch)
}
