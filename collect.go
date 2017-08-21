package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Describe - loops through the API metrics and passes them to prometheus.Describe
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {

	for _, m := range e.rancherMetrics {
		ch <- m
	}

}

// Collect function, called on by Prometheus Client library
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	eLogger.Debug("Metric collection requested")
	var data, err = e.gatherData()

	if err != nil {
		eLogger.Errorf("Error gathering data from metadata service, error: %s", err)
		return
	}

	if err := e.populateMetrics(data, ch); err != nil {
		eLogger.Printf("Error scraping rancher metadata service: %s", err)
		return
	}

	eLogger.Debug("Metrics successfully collected")

}

// populateMetrics - Collects the data from the API, returns data object
func (e *Exporter) populateMetrics(data *Data, ch chan<- prometheus.Metric) error {

	for _, x := range data.Services {

		// If system services have been ignored, the loop simply skips them
		if e.HideSys == true && x.System == true {
			continue
		}

		eLogger.Debug("Processing service metrics")

		// Checks the metric is of the expected kind
		if x.Kind != "service" {
			eLogger.Errorf("Error processing service metric, metric not of expected kind")
			continue
		}

		ch <- prometheus.MustNewConstMetric(e.rancherMetrics["servicesScale"], prometheus.GaugeValue, float64(x.Scale), x.Name, x.StackName)

		eLogger.Debug("Processing State metrics")

		for _, s := range e.ServiceStates {
			if x.State == s {
				ch <- prometheus.MustNewConstMetric(e.rancherMetrics["servicesState"], prometheus.GaugeValue, float64(1), x.Name, x.StackName, s)
			} else {
				ch <- prometheus.MustNewConstMetric(e.rancherMetrics["servicesState"], prometheus.GaugeValue, float64(0), x.Name, x.StackName, s)
			}

		}

	}

	for _, x := range data.Stacks {

		// If system stacks have been ignored, the loop simply skips them
		if e.HideSys == true && x.System == true {
			continue
		}

		eLogger.Debug("Processing stack metrics")

		// Checks the metric is of the expected kind
		if x.Kind != "stack" {
			eLogger.Errorf("Error processing stack metric, metric not of expected kind")
			continue
		}

		eLogger.Warn("Processing State metrics not currently possible, awaiting feature request in metadata service")

	}

	return nil
}
