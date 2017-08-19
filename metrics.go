package main

import "github.com/prometheus/client_golang/prometheus"

// addMetrics - Add's all of the GuageVecs to the `guageVecs` map, returns the map.
func addMetrics() map[string]*prometheus.Desc {

	gaugeVecs := make(map[string]*prometheus.Desc)

	// Service Metrics
	gaugeVecs["servicesScale"] = prometheus.NewDesc(
		prometheus.BuildFQName("rancher", "service", "scale"),
		"scale of defined service as reported by Rancher",
		[]string{"name", "stack_name"}, nil,
	)
	gaugeVecs["servicesHealth"] = prometheus.NewDesc(
		prometheus.BuildFQName("rancher", "service", "health_status"),
		"HealthState of the service, as reported by the Rancher API. Either (1) or (0)",
		[]string{"name", "stack_name", "health_state"}, nil,
	)
	gaugeVecs["servicesState"] = prometheus.NewDesc(
		prometheus.BuildFQName("rancher", "service", "state"),
		"State of the service, as reported by the Rancher API",
		[]string{"name", "stack_name", "state"}, nil,
	)

	// Stack Metrics

	// gaugeVecs["stacksHealth"] = prometheus.NewDesc(
	// 	prometheus.BuildFQName("rancher", "stack", "health_status"),
	// 	"HealthState of defined stack as reported by Rancher",
	// 	[]string{"name", "health_state", "system"}, nil,
	// )
	// gaugeVecs["stacksState"] = prometheus.NewDesc(
	// 	prometheus.BuildFQName("rancher", "stack", "health_status"),

	// 	prometheus.GaugeOpts{
	// 		Namespace: "rancher",
	// 		Name:      "stack_state",
	// 		Help:      "State of defined stack as reported by Rancher",
	// 	}, []string{"name", "state", "system"})

	// // Host Metrics
	// gaugeVecs["hostsState"] = prometheus.NewDesc(
	// 	prometheus.GaugeOpts{
	// 		Namespace: "rancher",
	// 		Name:      ("host_state"),
	// 		Help:      "State of defined host as reported by the Rancher API",
	// 	}, []string{"name", "state"})
	// gaugeVecs["hostAgentsState"] = prometheus.NewDesc(
	// 	prometheus.GaugeOpts{
	// 		Namespace: "rancher",
	// 		Name:      ("host_agent_state"),
	// 		Help:      "State of defined host agent as reported by the Rancher API",
	// 	}, []string{"name", "state"})

	return gaugeVecs
}
