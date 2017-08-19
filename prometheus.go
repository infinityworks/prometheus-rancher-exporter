package main

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// Describe - loops through the API metrics and passes them to prometheus.Describe
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {

	for _, m := range e.gaugeVecs {
		ch <- m
	}

}

// Collect function, called on by Prometheus Client library
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	log.Info("collecting")
	var data, err = e.gatherData()
	log.Info("collected")

	if err != nil {
		log.Error("Error getting JSON from URL ")
		return
	}

	if err := e.populateMetrics(data, ch); err != nil {
		log.Printf("Error scraping rancher metadata service: %s", err)
		return
	}

	log.Infof("Metrics successfully processed")

}
