# Example Metrics

Example of the metrics you could expect to see, returned for the service,stack and host states.

```
# HELP rancher_host_state State of defined host as reported by the Rancher API
# TYPE rancher_host_state gauge
rancher_host_state{name="example-server-01.c.rancher-dev.internal",state="activating"} 0
rancher_host_state{name="example-server-01.c.rancher-dev.internal",state="active"} 1
rancher_host_state{name="example-server-01.c.rancher-dev.internal",state="deactivating"} 0
rancher_host_state{name="example-server-01.c.rancher-dev.internal",state="error"} 0
rancher_host_state{name="example-server-01.c.rancher-dev.internal",state="erroring"} 0
rancher_host_state{name="example-server-01.c.rancher-dev.internal",state="inactive"} 0
rancher_host_state{name="example-server-01.c.rancher-dev.internal",state="provisioned"} 0
rancher_host_state{name="example-server-01.c.rancher-dev.internal",state="purged"} 0
rancher_host_state{name="example-server-01.c.rancher-dev.internal",state="purging"} 0
rancher_host_state{name="example-server-01.c.rancher-dev.internal",state="registering"} 0
rancher_host_state{name="example-server-01.c.rancher-dev.internal",state="removed"} 0
rancher_host_state{name="example-server-01.c.rancher-dev.internal",state="removing"} 0
rancher_host_state{name="example-server-01.c.rancher-dev.internal",state="requested"} 0
rancher_host_state{name="example-server-01.c.rancher-dev.internal",state="restoring"} 0
rancher_host_state{name="example-server-01.c.rancher-dev.internal",state="updating_active"} 0
rancher_host_state{name="example-server-01.c.rancher-dev.internal",state="updating_inactive"} 0
# HELP rancher_service_health_status HealthState of the service, as reported by the Rancher API. Either (1) or (0)
# TYPE rancher_service_health_status gauge
rancher_service_health_status{health_state="healthy",name="hubot",stack_name="rocket-chat"} 0
rancher_service_health_status{health_state="healthy",name="mongo",stack_name="rocket-chat"} 0
rancher_service_health_status{health_state="healthy",name="rocketchat",stack_name="rocket-chat"} 0
rancher_service_health_status{health_state="unhealthy",name="hubot",stack_name="rocket-chat"} 1
rancher_service_health_status{health_state="unhealthy",name="mongo",stack_name="rocket-chat"} 1
rancher_service_health_status{health_state="unhealthy",name="prometheus",stack_name="Prometheus"} 0
rancher_service_health_status{health_state="unhealthy",name="rocketchat",stack_name="rocket-chat"} 1
# HELP rancher_service_scale scale of defined service as reported by Rancher
# TYPE rancher_service_scale gauge
rancher_service_scale{name="hubot",stack_name="rocket-chat"} 1
rancher_service_scale{name="mongo",stack_name="rocket-chat"} 1
rancher_service_scale{name="rocketchat",stack_name="rocket-chat"} 1
# HELP rancher_service_state State of the service, as reported by the Rancher API
# TYPE rancher_service_state gauge
rancher_service_state{name="hubot",stack_name="rocket-chat",state="activating"} 0
rancher_service_state{name="hubot",stack_name="rocket-chat",state="active"} 0
rancher_service_state{name="hubot",stack_name="rocket-chat",state="canceled_upgrade"} 0
rancher_service_state{name="hubot",stack_name="rocket-chat",state="canceling_upgrade"} 0
rancher_service_state{name="hubot",stack_name="rocket-chat",state="deactivating"} 0
rancher_service_state{name="hubot",stack_name="rocket-chat",state="finishing_upgrade"} 0
rancher_service_state{name="hubot",stack_name="rocket-chat",state="inactive"} 1
rancher_service_state{name="hubot",stack_name="rocket-chat",state="registering"} 0
rancher_service_state{name="hubot",stack_name="rocket-chat",state="removed"} 0
rancher_service_state{name="hubot",stack_name="rocket-chat",state="removing"} 0
rancher_service_state{name="hubot",stack_name="rocket-chat",state="requested"} 0
rancher_service_state{name="hubot",stack_name="rocket-chat",state="restarting"} 0
rancher_service_state{name="hubot",stack_name="rocket-chat",state="rolling_back"} 0
rancher_service_state{name="hubot",stack_name="rocket-chat",state="updating_active"} 0
rancher_service_state{name="hubot",stack_name="rocket-chat",state="updating_inactive"} 0
rancher_service_state{name="hubot",stack_name="rocket-chat",state="upgraded"} 0
rancher_service_state{name="hubot",stack_name="rocket-chat",state="upgrading"} 0
rancher_service_state{name="mongo",stack_name="rocket-chat",state="activating"} 0
rancher_service_state{name="mongo",stack_name="rocket-chat",state="active"} 0
rancher_service_state{name="mongo",stack_name="rocket-chat",state="canceled_upgrade"} 0
rancher_service_state{name="mongo",stack_name="rocket-chat",state="canceling_upgrade"} 0
rancher_service_state{name="mongo",stack_name="rocket-chat",state="deactivating"} 0
rancher_service_state{name="mongo",stack_name="rocket-chat",state="finishing_upgrade"} 0
rancher_service_state{name="mongo",stack_name="rocket-chat",state="inactive"} 1
rancher_service_state{name="mongo",stack_name="rocket-chat",state="registering"} 0
rancher_service_state{name="mongo",stack_name="rocket-chat",state="removed"} 0
rancher_service_state{name="mongo",stack_name="rocket-chat",state="removing"} 0
rancher_service_state{name="mongo",stack_name="rocket-chat",state="requested"} 0
rancher_service_state{name="mongo",stack_name="rocket-chat",state="restarting"} 0
rancher_service_state{name="mongo",stack_name="rocket-chat",state="rolling_back"} 0
rancher_service_state{name="mongo",stack_name="rocket-chat",state="updating_active"} 0
rancher_service_state{name="mongo",stack_name="rocket-chat",state="updating_inactive"} 0
rancher_service_state{name="mongo",stack_name="rocket-chat",state="upgraded"} 0
rancher_service_state{name="mongo",stack_name="rocket-chat",state="upgrading"} 0
rancher_service_state{name="rocketchat",stack_name="rocket-chat",state="activating"} 0
rancher_service_state{name="rocketchat",stack_name="rocket-chat",state="active"} 0
rancher_service_state{name="rocketchat",stack_name="rocket-chat",state="canceled_upgrade"} 0
rancher_service_state{name="rocketchat",stack_name="rocket-chat",state="canceling_upgrade"} 0
rancher_service_state{name="rocketchat",stack_name="rocket-chat",state="deactivating"} 0
rancher_service_state{name="rocketchat",stack_name="rocket-chat",state="finishing_upgrade"} 0
rancher_service_state{name="rocketchat",stack_name="rocket-chat",state="inactive"} 1
rancher_service_state{name="rocketchat",stack_name="rocket-chat",state="registering"} 0
rancher_service_state{name="rocketchat",stack_name="rocket-chat",state="removed"} 0
rancher_service_state{name="rocketchat",stack_name="rocket-chat",state="removing"} 0
rancher_service_state{name="rocketchat",stack_name="rocket-chat",state="requested"} 0
rancher_service_state{name="rocketchat",stack_name="rocket-chat",state="restarting"} 0
rancher_service_state{name="rocketchat",stack_name="rocket-chat",state="rolling_back"} 0
rancher_service_state{name="rocketchat",stack_name="rocket-chat",state="updating_active"} 0
rancher_service_state{name="rocketchat",stack_name="rocket-chat",state="updating_inactive"} 0
rancher_service_state{name="rocketchat",stack_name="rocket-chat",state="upgraded"} 0
rancher_service_state{name="rocketchat",stack_name="rocket-chat",state="upgrading"} 0
# HELP rancher_stack_health_status HealthState of defined stack as reported by Rancher
# TYPE rancher_stack_health_status gauge
rancher_stack_health_status{health_state="healthy",name="rocket-chat"} 0
rancher_stack_health_status{health_state="unhealthy",name="rocket-chat"} 1
# HELP rancher_stack_state State of defined stack as reported by Rancher
# TYPE rancher_stack_state gauge
rancher_stack_state{name="rocket-chat",state="activating"} 0
rancher_stack_state{name="rocket-chat",state="active"} 1
rancher_stack_state{name="rocket-chat",state="canceled_upgrade"} 0
rancher_stack_state{name="rocket-chat",state="canceling_upgrade"} 0
rancher_stack_state{name="rocket-chat",state="error"} 0
rancher_stack_state{name="rocket-chat",state="erroring"} 0
rancher_stack_state{name="rocket-chat",state="finishing_upgrade"} 0
rancher_stack_state{name="rocket-chat",state="removed"} 0
rancher_stack_state{name="rocket-chat",state="removing"} 0
rancher_stack_state{name="rocket-chat",state="requested"} 0
rancher_stack_state{name="rocket-chat",state="restarting"} 0
rancher_stack_state{name="rocket-chat",state="rolling_back"} 0
rancher_stack_state{name="rocket-chat",state="updating_active"} 0
rancher_stack_state{name="rocket-chat",state="upgraded"} 0
rancher_stack_state{name="rocket-chat",state="upgrading"} 0
# HELP rancher_node_state State of defined node as reported by Rancher
# TYPE rancher_node_state gauge
rancher_node_state{cluster_name="cluster_name",node_name="node_name",state="active"} 1
rancher_node_state{cluster_name="cluster_name",node_name="node_name",state="cordoned"} 0
rancher_node_state{cluster_name="cluster_name",node_name="node_name",state="drained"} 0
rancher_node_state{cluster_name="cluster_name",node_name="node_name",state="draining"} 0
rancher_node_state{cluster_name="cluster_name",node_name="node_name",state="provisioning"} 0
rancher_node_state{cluster_name="cluster_name",node_name="node_name",state="registering"} 0
rancher_node_state{cluster_name="cluster_name",node_name="node_name",state="unavailable"} 0
# HELP rancher_cluster_component_status Component statuses of defined cluster as reported by Rancher
# TYPE rancher_cluster_component_status gauge
rancher_cluster_component_status{cluster_name="cluster_name",component_name="component_name",status="False"} 0
rancher_cluster_component_status{cluster_name="cluster_name",component_name="component_name",status="True"} 1
rancher_cluster_component_status{cluster_name="cluster_name",component_name="component_name",status="Unknown"} 0
# HELP rancher_cluster_state State of defined cluster as reported by Rancher
# TYPE gauge
rancher_cluster_state{cluster_name="cluster_name",state="active"} 1
rancher_cluster_state{cluster_name="cluster_name",state="cordoned"} 0
rancher_cluster_state{cluster_name="cluster_name",state="degraded"} 0
rancher_cluster_state{cluster_name="cluster_name",state="disconnected"} 0
rancher_cluster_state{cluster_name="cluster_name",state="drained"} 0
rancher_cluster_state{cluster_name="cluster_name",state="draining"} 0
rancher_cluster_state{cluster_name="cluster_name",state="healthy"} 0
rancher_cluster_state{cluster_name="cluster_name",state="initializing"} 0
rancher_cluster_state{cluster_name="cluster_name",state="locked"} 0
rancher_cluster_state{cluster_name="cluster_name",state="purged"} 0
rancher_cluster_state{cluster_name="cluster_name",state="purging"} 0
rancher_cluster_state{cluster_name="cluster_name",state="reconnecting"} 0
rancher_cluster_state{cluster_name="cluster_name",state="reinitializing"} 0
rancher_cluster_state{cluster_name="cluster_name",state="removed"} 0
rancher_cluster_state{cluster_name="cluster_name",state="running"} 0
rancher_cluster_state{cluster_name="cluster_name",state="unavailable"} 0
rancher_cluster_state{cluster_name="cluster_name",state="unhealthy"} 0
rancher_cluster_state{cluster_name="cluster_name",state="upgraded"} 0
rancher_cluster_state{cluster_name="cluster_name",state="upgrading"} 0
```

An example of the internal metrics to track the performance of the exporter, and useful as a basic example how to instrument your code.

```
# HELP function_count_totals total count of function calls
# TYPE function_count_totals counter
function_count_totals{fnc="getJSON",pkg="hosts"} 3
function_count_totals{fnc="getJSON",pkg="services"} 3
function_count_totals{fnc="getJSON",pkg="stacks"} 3
# HELP function_durations_seconds Function timings for Rancher Exporter
# TYPE function_durations_seconds summary
function_durations_seconds{fnc="getJSON",pkg="hosts",quantile="0.5"} 33546
function_durations_seconds{fnc="getJSON",pkg="hosts",quantile="0.9"} 59199
function_durations_seconds{fnc="getJSON",pkg="hosts",quantile="0.99"} 59199
function_durations_seconds_sum{fnc="getJSON",pkg="hosts"} 121502
function_durations_seconds_count{fnc="getJSON",pkg="hosts"} 3
function_durations_seconds{fnc="getJSON",pkg="services",quantile="0.5"} 49354
function_durations_seconds{fnc="getJSON",pkg="services",quantile="0.9"} 63310
function_durations_seconds{fnc="getJSON",pkg="services",quantile="0.99"} 63310
function_durations_seconds_sum{fnc="getJSON",pkg="services"} 146519
function_durations_seconds_count{fnc="getJSON",pkg="services"} 3
function_durations_seconds{fnc="getJSON",pkg="stacks",quantile="0.5"} 45075
function_durations_seconds{fnc="getJSON",pkg="stacks",quantile="0.9"} 59805
function_durations_seconds{fnc="getJSON",pkg="stacks",quantile="0.99"} 59805
function_durations_seconds_sum{fnc="getJSON",pkg="stacks"} 134789
function_durations_seconds_count{fnc="getJSON",pkg="stacks"} 3
```
