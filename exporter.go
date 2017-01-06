package main

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// Exporter Sets up all the runtime and metrics
type Exporter struct {
	rancherURL string
	accessKey  string
	secretKey  string
	hideSys    bool
	mutex      sync.RWMutex
	gaugeVecs  map[string]*prometheus.GaugeVec
}

// NewExporter creates the metrics we wish to monitor
func newExporter(rancherURL string, accessKey string, secretKey string, hideSys bool) *Exporter {

	gaugeVecs := addMetrics()
	return &Exporter{
		gaugeVecs:  gaugeVecs,
		rancherURL: rancherURL,
		accessKey:  accessKey,
		secretKey:  secretKey,
		hideSys:    hideSys,
	}
}
