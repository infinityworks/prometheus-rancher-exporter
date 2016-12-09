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
	rancherURL                 string
	accessKey                  string
	secretKey                  string
	mutex                      sync.RWMutex
	StackHealth                *prometheus.GaugeVec
	StackStateActivating       *prometheus.GaugeVec
	StackStateActive           *prometheus.GaugeVec
	StackStateCanceledUpgrade  *prometheus.GaugeVec
	StackStateCancelingUpgrade *prometheus.GaugeVec
	StackStateError            *prometheus.GaugeVec
	StackStateErroring         *prometheus.GaugeVec
	StackStateFinishingUpgrade *prometheus.GaugeVec
	StackStateRemoved          *prometheus.GaugeVec
	StackStateRemoving         *prometheus.GaugeVec
	StackStateRequested        *prometheus.GaugeVec
	StackStateRollingBack      *prometheus.GaugeVec
	StackStateUpdatingActive   *prometheus.GaugeVec
	StackStateUpgraded         *prometheus.GaugeVec
	StackStateUpgrading        *prometheus.GaugeVec
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
		StackStateActivating: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "stack_state_activating",
				Help:      "State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		StackStateActive: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "stack_state_active",
				Help:      "State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		StackStateCanceledUpgrade: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "stack_state_cancelled_upgrade",
				Help:      "State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		StackStateCancelingUpgrade: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "stack_state_canceling_upgrade",
				Help:      "State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		StackStateError: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "stack_state_error",
				Help:      "State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		StackStateErroring: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "stack_state_erroring",
				Help:      "State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		StackStateFinishingUpgrade: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "stack_state_finishing_upgrade",
				Help:      "State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		StackStateRemoved: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "stack_state_removed",
				Help:      "State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		StackStateRemoving: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "stack_state_removing",
				Help:      "State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		StackStateRequested: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "stack_state_requested",
				Help:      "State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		StackStateRollingBack: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "stack_state_rolling_back",
				Help:      "State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		StackStateUpdatingActive: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "stack_state_updating_active",
				Help:      "State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		StackStateUpgraded: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "stack_state_upgraded",
				Help:      "State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		StackStateUpgrading: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "stack_state_upgrading",
				Help:      "State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
	}

}

// Describe describes all the metrics ever exported by the Rancher exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.StackHealth.Describe(ch)
	e.StackStateActivating.Describe(ch)
	e.StackStateActive.Describe(ch)
	e.StackStateCanceledUpgrade.Describe(ch)
	e.StackStateCancelingUpgrade.Describe(ch)
	e.StackStateError.Describe(ch)
	e.StackStateErroring.Describe(ch)
	e.StackStateFinishingUpgrade.Describe(ch)
	e.StackStateRemoved.Describe(ch)
	e.StackStateRemoving.Describe(ch)
	e.StackStateRequested.Describe(ch)
	e.StackStateRollingBack.Describe(ch)
	e.StackStateUpdatingActive.Describe(ch)
	e.StackStateUpgraded.Describe(ch)
	e.StackStateUpgrading.Describe(ch)
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
	e.StackStateActivating.Reset()
	e.StackStateActive.Reset()
	e.StackStateCanceledUpgrade.Reset()
	e.StackStateCancelingUpgrade.Reset()
	e.StackStateError.Reset()
	e.StackStateErroring.Reset()
	e.StackStateFinishingUpgrade.Reset()
	e.StackStateRemoved.Reset()
	e.StackStateRemoving.Reset()
	e.StackStateRequested.Reset()
	e.StackStateRollingBack.Reset()
	e.StackStateUpdatingActive.Reset()
	e.StackStateUpgraded.Reset()
	e.StackStateUpgrading.Reset()

	var stacksEndpoint string

	if strings.Contains(rancherURL, "v1") {
		fmt.Println("Version 1 API detected, using legacy API fields")
		stacksEndpoint = "/environments/"

	} else if strings.Contains(rancherURL, "v2") {
		stacksEndpoint = "/stacls/"
	}

	fmt.Println("Scraping: ", rancherURL+stacksEndpoint)
	err, stacksData := getJSONstacks(rancherURL+stacksEndpoint, accessKey, secretKey)

	/*if strings.Contains(rancherURL, "v1") {
		fmt.Println("Version 1 API detected, using legacy API fields")
		fmt.Println("Scraping: ", rancherURL+"/environments/")
		err, stacksData = getJSONstacks(rancherURL+"/environments/", accessKey, secretKey)
	}
	else if strings.Contains(rancherURL, "v2") {
		fmt.Println("Scraping: ", rancherURL+"/stacks/")
		err, stacksData = getJSONstacks(rancherURL+"/stacks/", accessKey, secretKey)
	}
	*/

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

		var StackStateActivating float64
		var StackStateActive float64
		var StackStateCanceledUpgrade float64
		var StackStateCancelingUpgrade float64
		var StackStateError float64
		var StackStateErroring float64
		var StackStateFinishingUpgrade float64
		var StackStateRemoved float64
		var StackStateRemoving float64
		var StackStateRequested float64
		var StackStateRollingBack float64
		var StackStateUpdatingActive float64
		var StackStateUpgraded float64
		var StackStateUpgrading float64

		if x.State == "activating" {
			StackStateActivating = 1
		} else if x.State == "active" {
			StackStateActive = 1
		} else if x.State == "canceled-upgrade" {
			StackStateCanceledUpgrade = 1
		} else if x.State == "canceling-upgrade" {
			StackStateCancelingUpgrade = 1
		} else if x.State == "error" {
			StackStateError = 1
		} else if x.State == "erroring" {
			StackStateErroring = 1
		} else if x.State == "finishing-upgrade" {
			StackStateFinishingUpgrade = 1
		} else if x.State == "removed" {
			StackStateRemoved = 1
		} else if x.State == "removing" {
			StackStateRemoving = 1
		} else if x.State == "requested" {
			StackStateRequested = 1
		} else if x.State == "rolling-back" {
			StackStateRollingBack = 1
		} else if x.State == "updating-active" {
			StackStateUpdatingActive = 1
		} else if x.State == "upgraded" {
			StackStateUpgraded = 1
		} else if x.State == "upgrading" {
			StackStateUpgrading = 1
		}

		e.StackHealth.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackHealthState)
		e.StackStateActivating.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackStateActivating)
		e.StackStateActive.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackStateActive)
		e.StackStateCanceledUpgrade.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackStateCanceledUpgrade)
		e.StackStateCancelingUpgrade.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackStateCancelingUpgrade)
		e.StackStateError.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackStateError)
		e.StackStateErroring.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackStateErroring)
		e.StackStateFinishingUpgrade.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackStateFinishingUpgrade)
		e.StackStateRemoved.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackStateRemoved)
		e.StackStateRemoving.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackStateRemoving)
		e.StackStateRequested.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackStateRequested)
		e.StackStateRollingBack.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackStateRollingBack)
		e.StackStateUpdatingActive.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackStateUpdatingActive)
		e.StackStateUpgraded.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackStateUpgraded)
		e.StackStateUpgrading.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(StackStateUpgrading)
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
	e.StackStateActivating.Collect(ch)
	e.StackStateActive.Collect(ch)
	e.StackStateCanceledUpgrade.Collect(ch)
	e.StackStateCancelingUpgrade.Collect(ch)
	e.StackStateError.Collect(ch)
	e.StackStateErroring.Collect(ch)
	e.StackStateFinishingUpgrade.Collect(ch)
	e.StackStateRemoved.Collect(ch)
	e.StackStateRemoving.Collect(ch)
	e.StackStateRequested.Collect(ch)
	e.StackStateRollingBack.Collect(ch)
	e.StackStateUpdatingActive.Collect(ch)
	e.StackStateUpgraded.Collect(ch)
	e.StackStateUpgrading.Collect(ch)
}
