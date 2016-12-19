package hosts

import (
	"sync"

	"github.com/prometheus/log"

	"github.com/infinityworksltd/prometheus-rancher-exporter/utils"
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

func (e *Exporter) gatherMetrics(rancherURL string, accessKey string, secretKey string, ch chan<- prometheus.Metric) error {

	// Reset guageVecs back to 0
	for _, m := range e.gaugeVecs {
		m.Reset()
	}

	// Set the correct API endpoint for hosts
	endpoint := (rancherURL + "/hosts/")

	// Scrape EndPoint for JSON Data
	data := new(Data)
	err := utils.GetJson(endpoint, accessKey, secretKey, &data)

	if err != nil {
		log.Error("Error getting JSON from URL ", endpoint)
		return err
	}

	log.Info("JSON Fetched for hosts: ", data)

	// Host Metrics
	for _, x := range data.Data {

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
		log.Errorf("Error scraping rancher url: %s", err)
		return
	}

	for _, m := range e.gaugeVecs {
		m.Collect(ch)
	}
}
