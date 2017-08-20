package main

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// Describe - loops through the API metrics and passes them to prometheus.Describe
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {

	for _, m := range e.rancherMetrics {
		ch <- m
	}

}

// Collect function, called on by Prometheus Client library
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	log.Info("Metric collection requested")
	var data, err = e.gatherData()

	if err != nil {
		log.Errorf("Error gathering data from metadata service, error: %s", err)
		return
	}

	if err := e.populateMetrics(data, ch); err != nil {
		log.Printf("Error scraping rancher metadata service: %s", err)
		return
	}

	log.Infof("Metrics successfully processed")

}

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

		ch <- prometheus.MustNewConstMetric(e.rancherMetrics["servicesScale"], prometheus.GaugeValue, float64(x.Scale), x.Name, x.StackName)

		log.Info("Processing State metrics")
		for _, s := range e.ServiceStates {
			if x.State == s {
				log.Infof("Processing State metrics, %s", x.State, s)
				ch <- prometheus.MustNewConstMetric(e.rancherMetrics["servicesState"], prometheus.GaugeValue, float64(1), x.Name, x.StackName, s)
			} else {
				log.Infof("Processing State metrics, %s", x.State, s)
				ch <- prometheus.MustNewConstMetric(e.rancherMetrics["servicesState"], prometheus.GaugeValue, float64(0), x.Name, x.StackName, s)
			}

		}

	}

	for _, x := range data.Stacks {

		// If system stacks have been ignored, the loop simply skips them
		if e.HideSys == true && x.System == true {
			continue
		}

		log.Info("Processing stack metrics")

		// Checks the metric is of the expected kind
		if x.Kind != "stack" {
			log.Errorf("Error processing stack metric, metric not of expected kind")
			continue
		}

		log.Infof("Processing Stack metric. Stack: %s", x.Name)

		log.Info("Processing State metrics not currently possible, awaiting feature request in metadata service")

	}

	return nil
}
