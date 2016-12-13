package services

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

// ServicesData is used to store data from the services endpoint in the API
type ServicesData struct {
	Data []struct {
		HealthState string `json:"healthState"`
		Name        string `json:"name"`
		Scale       int    `json:"scale"`
		State       string `json:"state"`
	} `json:"data"`
}

//NewExporter creates the metrics we wish to monitor
func NewExporter(rancherURL string, accessKey string, secretKey string) *Exporter {

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
func getJSONservices(rancherURL string, accessKey string, secretKey string) (error, ServicesData) {
	pulledData := ServicesData{}
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

func (e *Exporter) serviceScrape(rancherURL string, accessKey string, secretKey string, ch chan<- prometheus.Metric) error {

	for _, m := range e.gaugeVecs {
		m.Reset()
	}

	fmt.Println("Scraping: ", rancherURL+"/services/")
	err, servicesData := getJSONservices(rancherURL+"/services/", accessKey, secretKey)
	if err != nil {
		return err
	}
	fmt.Println("JSON Fetched for services: ", servicesData)

	// Service Metrics
	for _, x := range servicesData.Data {

		var ServiceHealthState float64
		if x.HealthState == "healthy" {
			ServiceHealthState = 1
		}

		e.gaugeVecs["ServiceHealth"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceHealthState)
		e.gaugeVecs["ServiceScale"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(float64(x.Scale))

		// Set all the metrics to 0, unless we get a match
		e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "activating"}).Set(0)
		e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "active"}).Set(0)
		e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "canceled_upgrade"}).Set(0)
		e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "canceling_upgrade"}).Set(0)
		e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "deactivasting"}).Set(0)
		e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "finishing_upgrade"}).Set(0)
		e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "inactive"}).Set(0)
		e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "registering"}).Set(0)
		e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "removed"}).Set(0)
		e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "removing"}).Set(0)
		e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "requested"}).Set(0)
		e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "restarting"}).Set(0)
		e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "rolling_back"}).Set(0)
		e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "updating_active"}).Set(0)
		e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "updating_inactive"}).Set(0)
		e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "upgraded"}).Set(0)
		e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "upgrading"}).Set(0)

		// Match states of the API to known values and override our values above.
		if x.State == "activating" {
			e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "activating"}).Set(1)
		} else if x.State == "active" {
			e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "active"}).Set(1)
		} else if x.State == "canceled-upgrade" {
			e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "canceled_upgrade"}).Set(1)
		} else if x.State == "canceling-upgrade" {
			e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "canceling_upgrade"}).Set(1)
		} else if x.State == "deactivasting" {
			e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "deactivasting"}).Set(1)
		} else if x.State == "finishing-upgrade" {
			e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "finishing_upgrade"}).Set(1)
		} else if x.State == "inactive" {
			e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "inactive"}).Set(1)
		} else if x.State == "registering" {
			e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "registering"}).Set(1)
		} else if x.State == "removed" {
			e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "removed"}).Set(1)
		} else if x.State == "removing" {
			e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "removing"}).Set(1)
		} else if x.State == "requested" {
			e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "requested"}).Set(1)
		} else if x.State == "restarting" {
			e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "restarting"}).Set(1)
		} else if x.State == "rolling-back" {
			e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "rolling_back"}).Set(1)
		} else if x.State == "updating-active" {
			e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "updating_active"}).Set(1)
		} else if x.State == "updating-inactive" {
			e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "updating_inactive"}).Set(1)
		} else if x.State == "upgraded" {
			e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "upgraded"}).Set(1)
		} else if x.State == "upgrading" {
			e.gaugeVecs["ServiceState"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name, "state": "upgrading"}).Set(1)
		}

	}

	return nil
}

// Collect function, called on by Prometheus Client
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	if err := e.serviceScrape(e.rancherURL, e.accessKey, e.secretKey, ch); err != nil {
		log.Printf("Error scraping rancher url: %s", err)
		return
	}

	for _, m := range e.gaugeVecs {
		m.Collect(ch)
	}
}
