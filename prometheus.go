package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
)

// Resets the guageVecs back to 0
// Ensures we start from a clean sheet
func (e *Exporter) resetGaugeVecs() {

	for _, m := range e.gaugeVecs {
		m.Reset()
	}
}

// Describe describes all the metrics ever exported by the Rancher exporter
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {

	for _, m := range e.gaugeVecs {
		m.Describe(ch)
	}
}

// Collect function, called on by Prometheus Client library
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	e.resetGaugeVecs() // Clean starting point

	// Range over the pre-configured endpoints array
	for _, p := range endpoints {

		var data, err = e.gatherData(e.rancherURL, e.accessKey, e.secretKey, p, ch)

		if err != nil {
			log.Error("Error getting JSON from URL ", p)
			return
		}

		if err := e.processMetrics(data, p, e.hideSys, ch); err != nil {
			log.Printf("Error scraping rancher url: %s", err)
			return
		}

		log.Infof("Metrics successfully processed for %s", p)

	}

	for _, m := range e.gaugeVecs {
		m.Collect(ch)
	}

}
