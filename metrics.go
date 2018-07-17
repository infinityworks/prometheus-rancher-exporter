package main

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// addMetrics - Add's all of the GuageVecs to the `guageVecs` map, returns the map.
func addMetrics() map[string]*prometheus.GaugeVec {
	gaugeVecs := make(map[string]*prometheus.GaugeVec)

	// Stack Metrics
	gaugeVecs["stacksHealth"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "stack_health_status",
			Help:      "HealthState of defined stack as reported by Rancher",
		}, []string{"name", "health_state", "system"})
	gaugeVecs["stacksState"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "stack_state",
			Help:      "State of defined stack as reported by Rancher",
		}, []string{"name", "state", "system"})

	// Service Metrics
	gaugeVecs["servicesScale"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_scale",
			Help:      "scale of defined service as reported by Rancher",
		}, []string{"name", "stack_name"})
	gaugeVecs["servicesHealth"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_health_status",
			Help:      "HealthState of the service, as reported by the Rancher API. Either (1) or (0)",
		}, []string{"name", "stack_name", "health_state"})
	gaugeVecs["servicesState"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      "service_state",
			Help:      "State of the service, as reported by the Rancher API",
		}, []string{"name", "stack_name", "state"})

	// Host Metrics
	gaugeVecs["hostsState"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      ("host_state"),
			Help:      "State of defined host as reported by the Rancher API",
		}, []string{"name", "state"})
	gaugeVecs["hostAgentsState"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rancher",
			Name:      ("host_agent_state"),
			Help:      "State of defined host agent as reported by the Rancher API",
		}, []string{"name", "state"})

	return gaugeVecs
}

// checkMetric - Checks the base type stored in the API is correct, this ensures we are setting the right metric for the right endpoint.
func checkMetric(endpoint string, baseType string) bool {
	e := strings.TrimSuffix(endpoint, "s")

	// Backwards compatibility fix, the API in V1 wrong, this is to cover v1 usage.
	if baseType == "environment" && e == "stack" {
		return true
	} else if e == "service" && (baseType == "externalService" || baseType == "loadBalancerService") {
		return true
	} else if e != baseType {
		log.Errorf("API MisMatch, expected %s metric, got %s metric", e, baseType)
		return false
	}

	return true
}

// setServiceMetrics - Logic to set the state of a system as a gauge metric
func (e *Exporter) setServiceMetrics(name string, stack string, state string, health string, scale int) error {
	e.gaugeVecs["servicesScale"].With(prometheus.Labels{"name": name, "stack_name": stack}).Set(float64(scale))

	for _, y := range healthStates {
		if health == y {
			e.gaugeVecs["servicesHealth"].With(prometheus.Labels{"name": name, "stack_name": stack, "health_state": y}).Set(1)
		} else {
			e.gaugeVecs["servicesHealth"].With(prometheus.Labels{"name": name, "stack_name": stack, "health_state": y}).Set(0)
		}
	}

	for _, y := range serviceStates {
		if state == y {
			e.gaugeVecs["servicesState"].With(prometheus.Labels{"name": name, "stack_name": stack, "state": y}).Set(1)
		} else {
			e.gaugeVecs["servicesState"].With(prometheus.Labels{"name": name, "stack_name": stack, "state": y}).Set(0)
		}

	}
	return nil
}

// setStackMetrics - Logic to set the state of a system as a gauge metric
func (e *Exporter) setStackMetrics(name string, state string, health string, system string) error {
	for _, y := range healthStates {
		if health == y {
			e.gaugeVecs["stacksHealth"].With(prometheus.Labels{"name": name, "health_state": y, "system": system}).Set(1)
		} else {
			e.gaugeVecs["stacksHealth"].With(prometheus.Labels{"name": name, "health_state": y, "system": system}).Set(0)
		}
	}
	for _, y := range stackStates {
		if state == y {
			e.gaugeVecs["stacksState"].With(prometheus.Labels{"name": name, "state": y, "system": system}).Set(1)
		} else {
			e.gaugeVecs["stacksState"].With(prometheus.Labels{"name": name, "state": y, "system": system}).Set(0)
		}

	}
	return nil
}

// setHostMetrics - Logic to set the state of a system as a gauge metric
func (e *Exporter) setHostMetrics(name string, state, agentState string) error {
	for _, y := range hostStates {
		if state == y {
			e.gaugeVecs["hostsState"].With(prometheus.Labels{"name": name, "state": y}).Set(1)
		} else {
			e.gaugeVecs["hostsState"].With(prometheus.Labels{"name": name, "state": y}).Set(0)
		}

	}
	for _, y := range agentStates {
		if agentState == y {
			e.gaugeVecs["hostAgentsState"].With(prometheus.Labels{"name": name, "state": y}).Set(1)
		} else {
			e.gaugeVecs["hostAgentsState"].With(prometheus.Labels{"name": name, "state": y}).Set(0)
		}

	}
	return nil
}
