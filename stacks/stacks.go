package stacks

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
			Name:      "stack_health_state",
			Help:      "HealthState of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["StackStateActivating"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "stack_state_activating",
			Help:      "State of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["StackStateActive"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "stack_state_active",
			Help:      "State of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["StackStateCanceledUpgrade"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "stack_state_cancelled_upgrade",
			Help:      "State of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["StackStateCancelingUpgrade"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "stack_state_canceling_upgrade",
			Help:      "State of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["StackStateError"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "stack_state_error",
			Help:      "State of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["StackStateErroring"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "stack_state_erroring",
			Help:      "State of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["StackStateFinishingUpgrade"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "stack_state_finishing_upgrade",
			Help:      "State of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["StackStateRemoved"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "stack_state_removed",
			Help:      "State of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["StackStateRemoving"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "stack_state_removing",
			Help:      "State of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["StackStateRequested"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "stack_state_requested",
			Help:      "State of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["StackStateRollingBack"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "stack_state_rolling_back",
			Help:      "State of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["StackStateUpdatingActive"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "stack_state_updating_active",
			Help:      "State of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["StackStateUpgraded"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "stack_state_upgraded",
			Help:      "State of defined stack as reported by Rancher",
		}, []string{"rancherURL", "name"})
	gaugeVecs["StackStateUpgrading"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "stack_state_upgrading",
			Help:      "State of defined stack as reported by Rancher",
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

// Gets the JSON response from the API and places it in the stacksData
func getJSONstacks(collectionURL string, accessKey string, secretKey string) (error, StacksData) {

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
		e.gaugeVecs["StackStateActivating"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["StackStateActive"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["StackStateCanceledUpgrade"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["StackStateCancelingUpgrade"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["StackStateError"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["StackStateErroring"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["StackStateFinishingUpgrade"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["StackStateRemoved"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["StackStateRemoving"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["StackStateRequested"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["StackStateRollingBack"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["StackStateUpdatingActive"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["StackStateUpgraded"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)
		e.gaugeVecs["StackStateUpgrading"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(0)

		// Match states of the API to known values and override our values above.
		if x.State == "activating" {
			e.gaugeVecs["StackStateActivating"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "active" {
			e.gaugeVecs["StackStateActive"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "canceled-upgrade" {
			e.gaugeVecs["StackStateCanceledUpgrade"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "canceling-upgrade" {
			e.gaugeVecs["StackStateCancelingUpgrade"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "error" {
			e.gaugeVecs["StackStateError"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "erroring" {
			e.gaugeVecs["StackStateErroring"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "finishing-upgrade" {
			e.gaugeVecs["StackStateFinishingUpgrade"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "removed" {
			e.gaugeVecs["StackStateRemoved"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "removing" {
			e.gaugeVecs["StackStateRemoving"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "requested" {
			e.gaugeVecs["StackStateRequested"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "rolling-back" {
			e.gaugeVecs["StackStateRollingBack"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "updating-active" {
			e.gaugeVecs["StackStateUpdatingActive"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "upgraded" {
			e.gaugeVecs["StackStateUpgraded"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)
		} else if x.State == "upgrading" {
			e.gaugeVecs["StackStateUpgrading"].With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(1)

		}

	}
	return nil
}

// Collect function, called on by Prometheus Client
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

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
