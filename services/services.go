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
			Help:      "HealthState of defined service as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["ServiceStateActivating"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_state_activating",
			Help:      "Service State of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["ServiceStateActive"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_state_active",
			Help:      "Service State of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["ServiceStateCanceledUpgrade"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_state_canceled_upgrade",
			Help:      "HealthState of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["ServiceStateCancelingUpgrade"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_state_canceling_upgrade",
			Help:      "HealthState of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["ServiceStateDeactivating"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_state_deactivating",
			Help:      "HealthState of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["ServiceStateFinishingUpgrade"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_state_finishing_upgrade",
			Help:      "HealthState of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["ServiceStateInactive"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_state_inactive",
			Help:      "HealthState of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["ServiceStateRegistering"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_state_registering",
			Help:      "HealthState of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["ServiceStateRemoved"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_state_removed",
			Help:      "HealthState of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["ServiceStateRemoving"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_state_removing",
			Help:      "HealthState of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["ServiceStateRequested"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_state_requested",
			Help:      "HealthState of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["ServiceStateRestarting"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_state_restarting",
			Help:      "HealthState of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["ServiceStateRollingBack"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_state_rolling_back",
			Help:      "HealthState of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["ServiceStateUpdatingActive"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_state_updating_active",
			Help:      "HealthState of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["ServiceStateUpdatingInactive"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_state_updating_inactive",
			Help:      "HealthState of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["ServiceStateUpgraded"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_state_upgraded",
			Help:      "HealthState of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["ServiceStateUpgrading"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_state_upgrading",
			Help:      "HealthState of defined stack as reported by Rancher",
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
		e.gaugeVecs["ServiceStateActivating"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["ServiceStateActive"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["ServiceStateCanceledUpgrade"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["ServiceStateCancelingUpgrade"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["ServiceStateDeactivating"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["ServiceStateFinishingUpgrade"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["ServiceStateInactive"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["ServiceStateRegistering"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["ServiceStateRemoved"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["ServiceStateRemoving"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["ServiceStateRequested"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["ServiceStateRestarting"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["ServiceStateRollingBack"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["ServiceStateUpdatingActive"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["ServiceStateUpdatingInactive"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["ServiceStateUpgraded"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["ServiceStateUpgrading"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)

		// Match states of the API to known values and override our values above.
		if x.State == "activating" {
			e.gaugeVecs["ServiceStateActivating"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "active" {
			e.gaugeVecs["ServiceStateActive"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "canceled-upgrade" {
			e.gaugeVecs["ServiceStateCanceledUpgrade"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "canceling-upgrade" {
			e.gaugeVecs["ServiceStateCancelingUpgrade"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "deactivasting" {
			e.gaugeVecs["ServiceStateDeactivating"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "finishing-upgrade" {
			e.gaugeVecs["ServiceStateFinishingUpgrade"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "inactive" {
			e.gaugeVecs["ServiceStateInactive"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "registering" {
			e.gaugeVecs["ServiceStateRegistering"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "removed" {
			e.gaugeVecs["ServiceStateRemoved"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "removing" {
			e.gaugeVecs["ServiceStateRemoving"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "requested" {
			e.gaugeVecs["ServiceStateRequested"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "restarting" {
			e.gaugeVecs["ServiceStateRestarting"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "rolling-back" {
			e.gaugeVecs["ServiceStateRollingBack"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "updating-active" {
			e.gaugeVecs["ServiceStateUpdatingActive"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "updating-inactive" {
			e.gaugeVecs["ServiceStateUpdatingInactive"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "upgraded" {
			e.gaugeVecs["ServiceStateUpgraded"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "upgrading" {
			e.gaugeVecs["ServiceStateUpgrading"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
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
