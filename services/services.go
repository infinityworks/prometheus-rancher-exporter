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
	rancherURL                   string
	accessKey                    string
	secretKey                    string
	mutex                        sync.RWMutex
	ServiceScale                 *prometheus.GaugeVec
	ServiceHealth                *prometheus.GaugeVec
	ServiceStateActivating       *prometheus.GaugeVec
	ServiceStateActive           *prometheus.GaugeVec
	ServiceStateCanceledUpgrade  *prometheus.GaugeVec
	ServiceStateCancelingUpgrade *prometheus.GaugeVec
	ServiceStateDeactivating     *prometheus.GaugeVec
	ServiceStateFinishingUpgrade *prometheus.GaugeVec
	ServiceStateInactive         *prometheus.GaugeVec
	ServiceStateRegistering      *prometheus.GaugeVec
	ServiceStateRemoved          *prometheus.GaugeVec
	ServiceStateRemoving         *prometheus.GaugeVec
	ServiceStateRequested        *prometheus.GaugeVec
	ServiceStateRestarting       *prometheus.GaugeVec
	ServiceStateRollingBack      *prometheus.GaugeVec
	ServiceStateUpdatingActive   *prometheus.GaugeVec
	ServiceStateUpdatingInactive *prometheus.GaugeVec
	ServiceStateUpgraded         *prometheus.GaugeVec
	ServiceStateUpgrading        *prometheus.GaugeVec
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
				Name:      "service_health_status",
				Help:      "HealthState of defined service as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateActivating: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_activating",
				Help:      "Service State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateActive: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_active",
				Help:      "Service State of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateCanceledUpgrade: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_canceled_upgrade",
				Help:      "HealthState of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateCancelingUpgrade: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_canceling_upgrade",
				Help:      "HealthState of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateDeactivating: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_deactivating",
				Help:      "HealthState of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateFinishingUpgrade: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_finishing_upgrade",
				Help:      "HealthState of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateInactive: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_inactive",
				Help:      "HealthState of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateRegistering: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_registering",
				Help:      "HealthState of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateRemoved: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_removed",
				Help:      "HealthState of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateRemoving: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_removing",
				Help:      "HealthState of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateRequested: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_requested",
				Help:      "HealthState of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateRestarting: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_restarting",
				Help:      "HealthState of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateRollingBack: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_rolling_back",
				Help:      "HealthState of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateUpdatingActive: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_updating_active",
				Help:      "HealthState of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateUpdatingInactive: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_updating_inactive",
				Help:      "HealthState of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateUpgraded: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_upgraded",
				Help:      "HealthState of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
		ServiceStateUpgrading: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "rancher",
				Name:      "service_state_upgrading",
				Help:      "HealthState of defined stack as reported by Rancher",
			}, []string{"rancherURL", "name"}),
	}

}

// Describe describes all the metrics ever exported by the Rancher exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.ServiceScale.Describe(ch)
	e.ServiceHealth.Describe(ch)
	e.ServiceStateActivating.Describe(ch)
	e.ServiceStateActive.Describe(ch)
	e.ServiceStateCanceledUpgrade.Describe(ch)
	e.ServiceStateCancelingUpgrade.Describe(ch)
	e.ServiceStateDeactivating.Describe(ch)
	e.ServiceStateFinishingUpgrade.Describe(ch)
	e.ServiceStateInactive.Describe(ch)
	e.ServiceStateRegistering.Describe(ch)
	e.ServiceStateRemoved.Describe(ch)
	e.ServiceStateRemoving.Describe(ch)
	e.ServiceStateRequested.Describe(ch)
	e.ServiceStateRestarting.Describe(ch)
	e.ServiceStateRollingBack.Describe(ch)
	e.ServiceStateUpdatingActive.Describe(ch)
	e.ServiceStateUpdatingInactive.Describe(ch)
	e.ServiceStateUpgraded.Describe(ch)
	e.ServiceStateUpgrading.Describe(ch)
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
	e.ServiceStateActivating.Reset()
	e.ServiceStateActive.Reset()
	e.ServiceStateCanceledUpgrade.Reset()
	e.ServiceStateCancelingUpgrade.Reset()
	e.ServiceStateDeactivating.Reset()
	e.ServiceStateFinishingUpgrade.Reset()
	e.ServiceStateInactive.Reset()
	e.ServiceStateRegistering.Reset()
	e.ServiceStateRemoved.Reset()
	e.ServiceStateRemoving.Reset()
	e.ServiceStateRequested.Reset()
	e.ServiceStateRestarting.Reset()
	e.ServiceStateRollingBack.Reset()
	e.ServiceStateUpdatingActive.Reset()
	e.ServiceStateUpdatingInactive.Reset()
	e.ServiceStateUpgraded.Reset()
	e.ServiceStateUpgrading.Reset()

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

		var ServiceStateActivating float64
		var ServiceStateActive float64
		var ServiceStateCanceledUpgrade float64
		var ServiceStateCancelingUpgrade float64
		var ServiceStateDeactivating float64
		var ServiceStateFinishingUpgrade float64
		var ServiceStateInactive float64
		var ServiceStateRegistering float64
		var ServiceStateRemoved float64
		var ServiceStateRemoving float64
		var ServiceStateRequested float64
		var ServiceStateRestarting float64
		var ServiceStateRollingBack float64
		var ServiceStateUpdatingActive float64
		var ServiceStateUpdatingInactive float64
		var ServiceStateUpgraded float64
		var ServiceStateUpgrading float64

		if x.State == "activating" {
			ServiceStateActivating = 1
		} else if x.State == "active" {
			ServiceStateActive = 1
		} else if x.State == "canceled-upgrade" {
			ServiceStateCanceledUpgrade = 1
		} else if x.State == "canceling-upgrade" {
			ServiceStateCancelingUpgrade = 1
		} else if x.State == "deactivasting" {
			ServiceStateDeactivating = 1
		} else if x.State == "finishing-upgrade" {
			ServiceStateFinishingUpgrade = 1
		} else if x.State == "inactive" {
			ServiceStateInactive = 1
		} else if x.State == "registering" {
			ServiceStateRegistering = 1
		} else if x.State == "removed" {
			ServiceStateRemoved = 1
		} else if x.State == "removing" {
			ServiceStateRemoving = 1
		} else if x.State == "requested" {
			ServiceStateRequested = 1
		} else if x.State == "restarting" {
			ServiceStateRestarting = 1
		} else if x.State == "rolling-back" {
			ServiceStateRollingBack = 1
		} else if x.State == "updating-active" {
			ServiceStateUpdatingActive = 1
		} else if x.State == "updating-inactive" {
			ServiceStateUpdatingInactive = 1
		} else if x.State == "upgraded" {
			ServiceStateUpgraded = 1
		} else if x.State == "upgrading" {
			ServiceStateUpgrading = 1
		}

		e.ServiceHealth.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceHealthState)
		e.ServiceScale.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(float64(x.Scale))
		e.ServiceStateActivating.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateActivating)
		e.ServiceStateActive.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateActive)
		e.ServiceStateCanceledUpgrade.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateCanceledUpgrade)
		e.ServiceStateCancelingUpgrade.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateCancelingUpgrade)
		e.ServiceStateDeactivating.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateDeactivating)
		e.ServiceStateFinishingUpgrade.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateFinishingUpgrade)
		e.ServiceStateInactive.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateInactive)
		e.ServiceStateRegistering.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateRegistering)
		e.ServiceStateRemoved.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateRemoved)
		e.ServiceStateRemoving.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateRemoving)
		e.ServiceStateRequested.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateRequested)
		e.ServiceStateRestarting.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateRestarting)
		e.ServiceStateRollingBack.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateRollingBack)
		e.ServiceStateUpdatingActive.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateUpdatingActive)
		e.ServiceStateUpdatingInactive.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateUpdatingInactive)
		e.ServiceStateUpgraded.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateUpgraded)
		e.ServiceStateUpgrading.With(prometheus.Labels{"rancherURL": rancherURL, "name": x.Name}).Set(ServiceStateUpgrading)
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
	e.ServiceStateActivating.Collect(ch)
	e.ServiceStateActive.Collect(ch)
	e.ServiceStateCanceledUpgrade.Collect(ch)
	e.ServiceStateCancelingUpgrade.Collect(ch)
	e.ServiceStateDeactivating.Collect(ch)
	e.ServiceStateFinishingUpgrade.Collect(ch)
	e.ServiceStateInactive.Collect(ch)
	e.ServiceStateRegistering.Collect(ch)
	e.ServiceStateRemoved.Collect(ch)
	e.ServiceStateRemoving.Collect(ch)
	e.ServiceStateRequested.Collect(ch)
	e.ServiceStateRestarting.Collect(ch)
	e.ServiceStateRollingBack.Collect(ch)
	e.ServiceStateUpdatingActive.Collect(ch)
	e.ServiceStateUpdatingInactive.Collect(ch)
	e.ServiceStateUpgraded.Collect(ch)
	e.ServiceStateUpgrading.Collect(ch)
}
