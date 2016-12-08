package stacks

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
	rancherURL                  string
	accessKey                   string
	secretKey                   string
	mutex                       sync.RWMutex
	StackHealth                 *prometheus.GaugeVec
	StackStateHealthy           *prometheus.GaugeVec
	StackStateUnhealthy         *prometheus.GaugeVec
	StackStateUpdatingUnhealthy *prometheus.GaugeVec
	StackStateUpdatingHealthy   *prometheus.GaugeVec
	StackStateInitializing      *prometheus.GaugeVec
}

// StacksData is used to store data from the stacks endpoint in the API
type StacksData struct {
	Data []struct {
		HealthState string `json:"healthState"`
		Name        string `json:"name"`
		State       string `json:"state"`
	} `json:"data"`
}

func NewExporter(rancherURL string, accessKey string, secretKey string) *Exporter {
	return &Exporter{
		rancherURL: rancherURL,
		accessKey:  accessKey,
		secretKey:  secretKey,
		StackHealth: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "stack_health_state",
				Help:      "HealthState of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		StackStateHealthy: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "stack_state_healthy",
				Help:      "State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		StackStateUnhealthy: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "stack_state_unhealthy",
				Help:      "State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		StackStateUpdatingUnhealthy: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "stack_state_updating_unhealthy",
				Help:      "State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		StackStateUpdatingHealthy: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "stack_state_updating_healthy",
				Help:      "State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		StackStateInitializing: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "stack_state_initializing",
				Help:      "State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
	}

}

// Describe describes all the metrics ever exported by the Rancher exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.StackHealth.Describe(ch)
	e.StackStateHealthy.Describe(ch)
	e.StackStateUnhealthy.Describe(ch)
	e.StackStateUpdatingUnhealthy.Describe(ch)
	e.StackStateUpdatingHealthy.Describe(ch)
	e.StackStateInitializing.Describe(ch)
}

// Gets the JSON response from the API and places it in the struct
func getJSONstacks(rancherURL string, accessKey string, secretKey string) (error, StacksData) {
	pulledData := StacksData{}
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

func (e *Exporter) scrapeStacks(rancherURL string, accessKey string, secretKey string, ch chan<- prometheus.Metric) error {
	e.StackHealth.Reset()
	e.StackStateHealthy.Reset()
	e.StackStateUnhealthy.Reset()
	e.StackStateUpdatingUnhealthy.Reset()
	e.StackStateUpdatingHealthy.Reset()
	e.StackStateInitializing.Reset()

	fmt.Println("Scraping: ", rancherURL+"stacks/")
	err, stacksData := getJSONstacks(rancherURL+"stacks/", accessKey, secretKey)
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

		var StackStateHealthy float64
		var StackStateUnhealthy float64
		var StackStateUpdatingUnhealthy float64
		var StackStateUpdatingHealthy float64
		var StackStateInitializing float64

		if x.State == "healthy" {
			StackStateHealthy = 1
		}

		if x.State == "unhealthy" {
			StackStateUnhealthy = 1
		}

		if x.State == "updating-healthy" {
			StackStateUpdatingHealthy = 1
		}

		if x.State == "updating-unhealthy" {
			StackStateUpdatingUnhealthy = 1
		}

		if x.State == "initializing" {
			StackStateInitializing = 1
		}

		e.StackHealth.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackHealthState)
		e.StackStateHealthy.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackStateHealthy)
		e.StackStateUnhealthy.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackStateUnhealthy)
		e.StackStateUpdatingUnhealthy.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackStateUpdatingUnhealthy)
		e.StackStateUpdatingHealthy.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackStateUpdatingHealthy)
		e.StackStateInitializing.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackStateInitializing)
	}
	return nil
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	if err := e.scrapeStacks(e.rancherURL, e.accessKey, e.secretKey, ch); err != nil {
		log.Printf("Error scraping rancher url: %s", err)
		return
	}

	e.StackHealth.Collect(ch)
	e.StackStateHealthy.Collect(ch)
	e.StackStateUnhealthy.Collect(ch)
	e.StackStateUpdatingUnhealthy.Collect(ch)
	e.StackStateUpdatingHealthy.Collect(ch)
	e.StackStateInitializing.Collect(ch)
}
