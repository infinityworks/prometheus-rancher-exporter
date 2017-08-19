package main

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/infinityworks/prometheus-rancher-exporter/measure"
	"github.com/prometheus/client_golang/prometheus"
)

// Data is used to store data from all the relevant endpoints in the API
type Data struct {
	Services []struct {
		// HealthState string `json:"health_check"`
		Name      string `json:"name"`
		State     string `json:"state"`
		Scale     int    `json:"scale"`
		StackName string `json:"stack_name"`
		System    bool   `json:"system"`
		EnvID     string `json:"environment_uuid"`
		Kind      string `json:"kind"`
	} `json:"services"`
	Stacks []struct {
		Name string `json:"name"`
	}
	Hosts []struct {
		Name string `json:"name"`
	}
}

// type Data struct {
// 	Data []struct {
// 		HealthState string `json:"healthState"`
// 		Name        string `json:"name"`
// 		State       string `json:"state"`
// 		System      bool   `json:"system"`
// 		Scale       int    `json:"scale"`
// 		HostName    string `json:"hostname"`
// 		ID          string `json:"id"`
// 		StackID     string `json:"stackId"`
// 		EnvID       string `json:"environmentId"`
// 		BaseType    string `json:"basetype"`
// 		Type        string `json:"type"`
// 		AgentState  string `json:"agentState"`
// 	} `json:"data"`
// }

// populateMetrics - Collects the data from the API, returns data object
func (e *Exporter) populateMetrics(data *Data, ch chan<- prometheus.Metric) error {

	for _, x := range data.Services {

		// If system services have been ignored, the loop simply skips them
		if e.HideSys == true && x.System == true {
			continue
		}

		log.Info("Processing service metrics")

		// Checks the metric is of the expected kind
		if x.Kind != "service" {
			log.Errorf("Error processing service metric, metric not of expected kind")
			continue
		}

		log.Infof("Processing Scale metric. Service: %s, Stack %s, Scale %v", x.Name, x.StackName, float64(x.Scale))

		ch <- prometheus.MustNewConstMetric(e.gaugeVecs["servicesScale"], prometheus.GaugeValue, float64(x.Scale), x.Name, x.StackName)

		log.Info("Processing State metrics")
		for _, s := range e.ServiceStates {
			if x.State == s {
				log.Infof("Processing State metrics, %s", x.State, s)
				ch <- prometheus.MustNewConstMetric(e.gaugeVecs["servicesState"], prometheus.GaugeValue, float64(1), x.Name, x.StackName, s)
			} else {
				log.Infof("Processing State metrics, %s", x.State, s)
				ch <- prometheus.MustNewConstMetric(e.gaugeVecs["servicesState"], prometheus.GaugeValue, float64(0), x.Name, x.StackName, s)
			}

		}

	}

	return nil
}

// gatherData - Collects the data from thw API, invokes functions to transform that data into metrics
func (e *Exporter) gatherData() (*Data, error) {

	// Create new data slice from Struct
	var data = new(Data)

	// Scrape EndPoint for JSON Data
	err := getJSON(e.MetaDataURL, &data)
	if err != nil {
		log.Errorf("Error getting JSON from URL: %s", err)
		return nil, err
	}
	log.Debugf("JSON Fetched : ", data)

	return data, err
}

// getJSON return json from server, return the formatted JSON
func getJSON(url string, target interface{}) error {

	start := time.Now()

	// Counter for internal exporter metrics
	measure.FunctionCountTotal.With(prometheus.Labels{"pkg": "main", "fnc": "getJSON"}).Inc()

	log.Info("Scraping: ", url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")

	log.Info("header set")

	resp, err := client.Do(req)

	if err != nil {
		log.Error("Error Collecting JSON from API: ", err)
		panic(err)
	}

	respFormatted := json.NewDecoder(resp.Body).Decode(target)
	log.Info("formatted")

	// Timings recorded as part of internal metrics
	elapsed := float64((time.Since(start)) / time.Microsecond)
	measure.FunctionDurations.WithLabelValues("main", "getJSON").Observe(elapsed)

	// Close the response body, the underlying Transport should then close the connection.
	resp.Body.Close()

	// return formatted JSON
	return respFormatted
}
