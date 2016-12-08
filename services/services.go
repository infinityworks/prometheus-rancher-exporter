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
	rancherURL                    string
	accessKey                     string
	secretKey                     string
	mutex                         sync.RWMutex
	ServiceScale                  *prometheus.GaugeVec
	ServiceHealth                 *prometheus.GaugeVec
	ServiceStateHealthy           *prometheus.GaugeVec
	ServiceStateUnhealthy         *prometheus.GaugeVec
	ServiceStateUpdatingUnhealthy *prometheus.GaugeVec
	ServiceStateUpdatingHealthy   *prometheus.GaugeVec
	ServiceStateInitializing      *prometheus.GaugeVec
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

// NewExporter returns an initialized Exporter.
func NewExporter(rancherURL string, accessKey string, secretKey string) *Exporter {
	return &Exporter{
		rancherURL: rancherURL,
		accessKey:  accessKey,
		secretKey:  secretKey,
		ServiceScale: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_scale",
				Help:      "scale of defined service as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceHealth: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_health_state",
				Help:      "HealthState of defined service as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateHealthy: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_healthy",
				Help:      "Service State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateUnhealthy: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_unhealthy",
				Help:      "Service State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateUpdatingUnhealthy: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_updating_unhealthy",
				Help:      "HealthState of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateUpdatingHealthy: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_updating_healthy",
				Help:      "HealthState of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateInitializing: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_initializing",
				Help:      "HealthState of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
	}

}

// Describe describes all the metrics ever exported by the Rancher exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.ServiceScale.Describe(ch)
	e.ServiceHealth.Describe(ch)
	e.ServiceStateHealthy.Describe(ch)
	e.ServiceStateUnhealthy.Describe(ch)
	e.ServiceStateUpdatingUnhealthy.Describe(ch)
	e.ServiceStateUpdatingHealthy.Describe(ch)
	e.ServiceStateInitializing.Describe(ch)
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
	e.ServiceScale.Reset()
	e.ServiceHealth.Reset()
	e.ServiceStateHealthy.Reset()
	e.ServiceStateUnhealthy.Reset()
	e.ServiceStateUpdatingUnhealthy.Reset()
	e.ServiceStateUpdatingHealthy.Reset()
	e.ServiceStateInitializing.Reset()

	fmt.Println("Scraping: ", rancherURL+"services/")
	//err, servicesData := getJSONservices(rancherURL+"services/", accessKey, secretKey)
	err, servicesData := getJSONservices(rancherURL+"services/", accessKey, secretKey)
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

		var ServiceStateHealthy float64
		var ServiceStateUnhealthy float64
		var ServiceStateUpdatingUnhealthy float64
		var ServiceStateUpdatingHealthy float64
		var ServiceStateInitializing float64

		if x.State == "healthy" {
			ServiceStateHealthy = 1

		}

		if x.State == "unhealthy" {
			ServiceStateUnhealthy = 1
		}

		if x.State == "updating-healthy" {
			ServiceStateUpdatingHealthy = 1
		}

		if x.State == "updating-unhealthy" {
			ServiceStateUpdatingUnhealthy = 1
		}

		if x.State == "initializing" {
			ServiceStateInitializing = 1
		}

		e.ServiceHealth.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceHealthState)
		e.ServiceScale.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(float64(x.Scale))
		e.ServiceHealth.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceHealthState)
		e.ServiceStateHealthy.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateHealthy)
		e.ServiceStateUnhealthy.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateUnhealthy)
		e.ServiceStateUpdatingUnhealthy.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateUpdatingUnhealthy)
		e.ServiceStateUpdatingHealthy.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateUpdatingHealthy)
		e.ServiceStateInitializing.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateInitializing)

	}

	return nil
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	if err := e.serviceScrape(e.rancherURL, e.accessKey, e.secretKey, ch); err != nil {
		log.Printf("Error scraping rancher url: %s", err)
		return
	}

	e.ServiceScale.Collect(ch)
	e.ServiceHealth.Collect(ch)
	e.ServiceStateHealthy.Collect(ch)
	e.ServiceStateUnhealthy.Collect(ch)
	e.ServiceStateUpdatingUnhealthy.Collect(ch)
	e.ServiceStateUpdatingHealthy.Collect(ch)
	e.ServiceStateInitializing.Collect(ch)
}
