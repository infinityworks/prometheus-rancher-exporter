package stacks

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
	hideSys    bool
	mutex      sync.RWMutex
	gaugeVecs  map[string]*prometheus.GaugeVec
}

// Data is used to store data from the stacks endpoint in the API
type Data struct {
	Data []struct {
		HealthState string `json:"healthState"`
		Name        string `json:"name"`
		State       string `json:"state"`
		System      bool   `json:"system"`
	} `json:"data"`
}

//NewExporter creates the metrics we wish to monitor
func NewExporter(rancherURL string, accessKey string, secretKey string, hideSys bool) *Exporter {

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

func (e *Exporter) gatherMetrics(rancherURL string, accessKey string, secretKey string, hideSys bool, ch chan<- prometheus.Metric) error {

	// Reset guageVecs back to 0
	for _, m := range e.gaugeVecs {
		m.Reset()
	}

	// Check API version and return the correct URL path
	endpoint := utils.StacksURLCheck(rancherURL)

	// Scrape EndPoint for JSON Data
	data := new(Data)
	err := utils.GetJson(endpoint, accessKey, secretKey, &data)

	if err != nil {
		log.Error("Error getting JSON from URL ", endpoint)
		return err
	}

	log.Info("JSON Fetched for stacks: ", data)

	// Stack Metrics
	for _, x := range data.Data {

		// If system services have been ignored, the loop simply skips them
		if hideSys == true && x.System == true {
			continue
		} else {

			// Get the healthy state for a stack
			if x.HealthState == "healthy" {
				e.gaugeVecs["StackHealth"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
			} else {
				e.gaugeVecs["StackHealth"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
			}

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
	}
	return nil
}

// Collect function, called on by Prometheus Client
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	if err := e.gatherMetrics(e.rancherURL, e.accessKey, e.secretKey, e.hideSys, ch); err != nil {
		log.Errorf("Error scraping rancher url: %s", err)
		return
	}
	for _, m := range e.gaugeVecs {
		m.Collect(ch)
	}

}
