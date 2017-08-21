package metrics

import "github.com/prometheus/client_golang/prometheus"

// Return - returns a map of metrics to be used by the exporter
func Return() map[string]*prometheus.Desc {

	rancherMetrics := make(map[string]*prometheus.Desc)

	// Rancher Service Metrics
	rancherMetrics["servicesScale"] = prometheus.NewDesc(
		prometheus.BuildFQName("rancher", "service", "scale"),
		"scale of defined service as reported by Rancher",
		[]string{"name", "stack_name"}, nil,
	)
	rancherMetrics["servicesHealth"] = prometheus.NewDesc(
		prometheus.BuildFQName("rancher", "service", "health_status"),
		"HealthState of the service, as reported by the Rancher API. Either (1) or (0)",
		[]string{"name", "stack_name", "health_state"}, nil,
	)
	rancherMetrics["servicesState"] = prometheus.NewDesc(
		prometheus.BuildFQName("rancher", "service", "state"),
		"State of the service, as reported by the Rancher API",
		[]string{"name", "stack_name", "state"}, nil,
	)

	// Rancher Stack Metrics
	rancherMetrics["stacksHealth"] = prometheus.NewDesc(
		prometheus.BuildFQName("rancher", "stack", "health_status"),
		"HealthState of defined stack as reported by Rancher",
		[]string{"name", "health_state", "system"}, nil,
	)
	rancherMetrics["stacksState"] = prometheus.NewDesc(
		prometheus.BuildFQName("rancher", "stack", "state"),
		"State of defined stack as reported by Rancher",
		[]string{"name", "state", "system"}, nil,
	)

	// Rancher Host Metrics
	// rancherMetrics["hostsState"] = prometheus.NewDesc(
	// prometheus.BuildFQName("rancher", "host", "state"),
	// 	"State of defined host as reported by the Rancher API",
	// 	[]string{"name", "docker_version", "state"}, nil,
	// )
	// rancherMetrics["hostAgentsState"] = prometheus.NewDesc(
	// prometheus.BuildFQName("rancher", "host", "agent_state"),
	//	"State of defined host agent as reported by the Rancher API",
	// 	[]string{"name", "state"}, nil,
	// )

	return rancherMetrics
}

//
