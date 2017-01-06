package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/infinityworksltd/prometheus-rancher-exporter/measure"
	"github.com/prometheus/client_golang/prometheus"
)

// Data is used to store data from all the relevant endpoints in the API
type Data struct {
	Data []struct {
		HealthState string `json:"healthState"`
		Name        string `json:"name"`
		State       string `json:"state"`
		System      bool   `json:"system"`
		Scale       int    `json:"scale"`
		HostName    string `json:"hostname"`
		ID          string `json:"id"`
		StackID     string `json:"stackId"`
		EnvID       string `json:"environmentId"`
		BaseType    string `json:"basetype"`
	} `json:"data"`
}

// processMetrics - Collects the data from the API, returns data object
func (e *Exporter) processMetrics(data *Data, endpoint string, hideSys bool, ch chan<- prometheus.Metric) error {

	// Used for backwards compatibility
	apiVer := getAPIVersion(rancherURL)

	// Metrics - range through the data object
	for _, x := range data.Data {

		// If system services have been ignored, the loop simply skips them
		if hideSys == true && x.System == true {
			continue
		}

		// Checks the metric is of the expected type
		if checkMetric(endpoint, x.BaseType) == false {
			continue
		}

		log.Debug("Processing metrics for %s", endpoint)

		if endpoint == "hosts" {

			if err := e.setHostMetrics(x.HostName, x.State); err != nil {
				log.Errorf("Error processing host metrics: %s", err)
				log.Errorf("Attempt Failed to set %s, %s ", x.HostName, x.State)

				continue
			}

		} else if endpoint == "stacks" {

			// Used to create a map of stackID and stackName
			// Later used as a dimension in service metrics
			stackRef = storeStackRef(x.ID, x.Name)

			if err := e.setStackMetrics(x.Name, x.State, x.HealthState); err != nil {
				log.Errorf("Error processing stack metrics: %s", err)
				log.Errorf("Attempt Failed to set %s, %s, %s ", x.Name, x.State, x.HealthState)
				continue
			}

		} else if endpoint == "services" {

			// Retrieves the stack Name from the previous values stored.
			var stackName = ""
			if apiVer == "v1" {
				stackName = retrieveStackRef(x.EnvID)
			} else {
				stackName = retrieveStackRef(x.StackID)
			}

			if stackName == "unknown" {
				log.Warnf("Failed to obtain stack_name for %s from the API", x.Name)
			}

			if err := e.setServiceMetrics(x.Name, stackName, x.State, x.HealthState, x.Scale); err != nil {
				log.Errorf("Error processing service metrics: %s", err)
				log.Errorf("Attempt Failed to set %s, %s, %s, %s, %d", x.Name, stackName, x.State, x.HealthState, x.Scale)
				continue
			}

			e.setServiceMetrics(x.Name, stackName, x.State, x.HealthState, x.Scale)
		}

	}

	return nil
}

// gatherData - Collects the data from thw API, invokes functions to transform that data into metrics
func (e *Exporter) gatherData(rancherURL string, accessKey string, secretKey string, endpoint string, ch chan<- prometheus.Metric) (*Data, error) {

	// Check API version and return the correct URL path
	apiVer := getAPIVersion(rancherURL)
	url := setEndpoint(rancherURL, endpoint, apiVer)

	// Create new data slice from Struct
	var data = new(Data)

	// Scrape EndPoint for JSON Data
	err := getJSON(url, accessKey, secretKey, &data)
	if err != nil {
		log.Error("Error getting JSON from URL ", endpoint)
		return nil, err
	}
	log.Debugf("JSON Fetched for: "+endpoint+": ", data)

	return data, err
}

// getJSON return json from server, return the formatted JSON
func getJSON(url string, accessKey string, secretKey string, target interface{}) error {

	start := time.Now()

	// Counter for internal exporter metrics
	measure.FunctionCountTotal.With(prometheus.Labels{"pkg": "main", "fnc": "getJSON"}).Inc()

	log.Info("Scraping: ", url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(accessKey, secretKey)
	resp, err := client.Do(req)

	if err != nil {
		log.Error("Error Collecting JSON from API: ", err)
		panic(err)
	}

	respFormatted := json.NewDecoder(resp.Body).Decode(target)

	// Timings recorded as part of internal metrics
	elapsed := float64((time.Since(start)) / time.Microsecond)
	measure.FunctionDurations.WithLabelValues("hosts", "getJSON").Observe(elapsed)

	// Close the response body, the underlying Transport should then close the connection.
	resp.Body.Close()

	// return formatted JSON
	return respFormatted
}

// setEndpoint - Determines the correct URL endpoint to use, gives us backwards compatibility
func setEndpoint(rancherURL string, component string, apiVer string) string {

	var endpoint string

	if strings.Contains(component, "services") {
		endpoint = (rancherURL + "/services/")
	} else if strings.Contains(component, "hosts") {
		endpoint = (rancherURL + "/hosts/")
	} else if strings.Contains(component, "stacks") {

		if apiVer == "v1" {
			endpoint = (rancherURL + "/environments/")
		} else {
			endpoint = (rancherURL + "/stacks/")
		}
	}

	return endpoint
}

// getAPIVersion - Determines the API version in-use
func getAPIVersion(rancherURL string) string {

	var apiVer string

	if strings.Contains(rancherURL, "v1") {
		log.Debug("Version 1 API detected, using legacy API fields")
		apiVer = ("v1")

	} else if strings.Contains(rancherURL, "v2") {
		log.Debug("Version 2 API detected, using legacy API fields")
		apiVer = ("v2")

	} else {
		log.Info("Unknown API version detected, defaulting to v2")
		apiVer = ("v2")
	}

	return apiVer
}

// storeStackRef stores the stackID and stack name for use as a label elsewhere
func storeStackRef(stackID string, stackName string) map[string]string {

	stackRef[stackID] = stackName

	return stackRef
}

// retrieveStackRef returns the stack name, when sending the stackID
func retrieveStackRef(stackID string) string {

	for key, value := range stackRef {
		if stackID == "" {
			return "unknown"
		} else if stackID == key {
			log.Debugf("StackRef - Key is %s, Value is %s StackID is %s", key, value, stackID)
			return value
		}
	}
	// returns unknown if no match was found
	return "unknown"
}
