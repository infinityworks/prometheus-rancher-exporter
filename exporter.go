package main

import (
	"regexp"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// Exporter Sets up all the runtime and metrics
type Exporter struct {
	labelsFilter *regexp.Regexp
	rancherURL   string
	accessKey    string
	secretKey    string
	hideSys      bool
	mutex        sync.RWMutex
	gaugeVecs    map[string]*prometheus.GaugeVec
}

// NewExporter creates the metrics we wish to monitor
func newExporter(rancherURL, accessKey, secretKey string, labelsFilter *regexp.Regexp, hideSys bool) *Exporter {
	gaugeVecs := addMetrics()
	return &Exporter{
		labelsFilter: labelsFilter,
		gaugeVecs:    gaugeVecs,
		rancherURL:   rancherURL,
		accessKey:    accessKey,
		secretKey:    secretKey,
		hideSys:      hideSys,
	}
}
