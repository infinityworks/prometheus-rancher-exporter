# prometheus-rancher-exporter

Exposes the health of stacks/services and hosts from the Rancher API, to a Prometheus compatible endpoint.

Please Note this exporter has been re-written for the new v2 Rancher API available from version 1.2 of Rancher onwards. 
Testing has focused on the v1/v2-beta API's available on Rancher 1.2, If you require a version compatible with older versions, please use versions <05 from the Dockerhub.

## Description

This container can make use of Ranchers ability to assign API access to a container at runtime. This is achieved through labels to create a connection to the API.
Details of how to set this are shown below in the example rancher-compose configuration. 

The application, expects to get the following environment variables from the host, if you are using this externally to Rancher, or without the use of the labels to obtain an API key, you can update these values yourself, using environment variables.

Required:
* CATTLE_ACCESS_KEY
* CATTLE_SECRET_KEY
* CATTLE_URL

Optional
* METRICS_PATH  //Path under which to expose metrics.
* LISTEN_ADDRESS // Port on which to expose metrics.

## Install and deploy

Run manually from Docker Hub:
```
docker run -d -e CATTLE_ACCESS_KEY="XXXXXXXX" -e CATTLE_SECRET_KEY="XXXXXXX" -e CATTLE_URL="http://<YOUR_IP>:8080/v2-beta" -p 9010:9010 infinityworks/prometheus-rancher-exporter
```

Build a docker image:
```
docker build -t <image-name> .
docker run -d -e CATTLE_ACCESS_KEY="XXXXXXXX" -e CATTLE_SECRET_KEY="XXXXXXX" -e CATTLE_URL="http://<YOUR_IP>:8080/v2-beta" -p 9010:9010 <image-name>
```

Running the node process:
```
DEBUG=re node app.js
```

## Docker compose

```
prometheus-rancher-exporter:
    tty: true
    stdin_open: true
    labels:
      io.rancher.container.create_agent: true
      io.rancher.container.agent.role: environment
    expose:
      - 9010:9010
    image: infinityworks/prometheus-rancher-exporter:latest
```

## Metrics

Metrics will be made available on port 9010 by default, or you can pass environment variable ```LISTEN_ADDRESS``` to override this.

Below are all the metrics currently exporterd. The metrics are broken down by host/stack/service and map out to the individual states available in the Rancher API.

```
# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
go_gc_duration_seconds{quantile="0.5"} 0
go_gc_duration_seconds{quantile="0.75"} 0
go_gc_duration_seconds{quantile="1"} 0
go_gc_duration_seconds_sum 0
go_gc_duration_seconds_count 0
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 9
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 765008
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 765008
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 1.442629e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 213
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 65536
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 765008
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 221184
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 1.482752e+06
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 6541
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 0
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 1.703936e+06
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 0
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 18
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 6754
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 4800
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 16384
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 18240
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 32768
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 4.194304e+06
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 1.066419e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 393216
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 393216
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 4.720888e+06
# HELP http_request_duration_microseconds The HTTP request latencies in microseconds.
# TYPE http_request_duration_microseconds summary
http_request_duration_microseconds{handler="prometheus",quantile="0.5"} NaN
http_request_duration_microseconds{handler="prometheus",quantile="0.9"} NaN
http_request_duration_microseconds{handler="prometheus",quantile="0.99"} NaN
http_request_duration_microseconds_sum{handler="prometheus"} 0
http_request_duration_microseconds_count{handler="prometheus"} 0
# HELP http_request_size_bytes The HTTP request sizes in bytes.
# TYPE http_request_size_bytes summary
http_request_size_bytes{handler="prometheus",quantile="0.5"} NaN
http_request_size_bytes{handler="prometheus",quantile="0.9"} NaN
http_request_size_bytes{handler="prometheus",quantile="0.99"} NaN
http_request_size_bytes_sum{handler="prometheus"} 0
http_request_size_bytes_count{handler="prometheus"} 0
# HELP http_response_size_bytes The HTTP response sizes in bytes.
# TYPE http_response_size_bytes summary
http_response_size_bytes{handler="prometheus",quantile="0.5"} NaN
http_response_size_bytes{handler="prometheus",quantile="0.9"} NaN
http_response_size_bytes{handler="prometheus",quantile="0.99"} NaN
http_response_size_bytes_sum{handler="prometheus"} 0
http_response_size_bytes_count{handler="prometheus"} 0
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 0.03
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1.048576e+06
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 8
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 7.213056e+06
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.48128981539e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 1.1776e+07
# HELP rancher_host_state_activating State of defined host as reported by Rancher
# TYPE rancher_host_state_activating gauge
rancher_host_state_activating{name="example-server.rancher-dev.internal",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_host_state_active State of defined host as reported by Rancher
# TYPE rancher_host_state_active gauge
rancher_host_state_active{name="example-server.rancher-dev.internal",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
# HELP rancher_host_state_deactivating State of defined host as reported by Rancher
# TYPE rancher_host_state_deactivating gauge
rancher_host_state_deactivating{name="example-server.rancher-dev.internal",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_host_state_error State of defined host as reported by Rancher
# TYPE rancher_host_state_error gauge
rancher_host_state_error{name="example-server.rancher-dev.internal",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_host_state_erroring State of defined host as reported by Rancher
# TYPE rancher_host_state_erroring gauge
rancher_host_state_erroring{name="example-server.rancher-dev.internal",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_host_state_inactive State of defined host as reported by Rancher
# TYPE rancher_host_state_inactive gauge
rancher_host_state_inactive{name="example-server.rancher-dev.internal",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_host_state_provisioned State of defined host as reported by Rancher
# TYPE rancher_host_state_provisioned gauge
rancher_host_state_provisioned{name="example-server.rancher-dev.internal",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_host_state_purged State of defined host as reported by Rancher
# TYPE rancher_host_state_purged gauge
rancher_host_state_purged{name="example-server.rancher-dev.internal",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_host_state_purging State of defined host as reported by Rancher
# TYPE rancher_host_state_purging gauge
rancher_host_state_purging{name="example-server.rancher-dev.internal",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_host_state_registering State of defined host as reported by Rancher
# TYPE rancher_host_state_registering gauge
rancher_host_state_registering{name="example-server.rancher-dev.internal",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_host_state_removed State of defined host as reported by Rancher
# TYPE rancher_host_state_removed gauge
rancher_host_state_removed{name="example-server.rancher-dev.internal",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_host_state_removing State of defined host as reported by Rancher
# TYPE rancher_host_state_removing gauge
rancher_host_state_removing{name="example-server.rancher-dev.internal",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_host_state_requested State of defined host as reported by Rancher
# TYPE rancher_host_state_requested gauge
rancher_host_state_requested{name="example-server.rancher-dev.internal",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_host_state_restoring State of defined host as reported by Rancher
# TYPE rancher_host_state_restoring gauge
rancher_host_state_restoring{name="example-server.rancher-dev.internal",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_host_state_updating_active State of defined host as reported by Rancher
# TYPE rancher_host_state_updating_active gauge
rancher_host_state_updating_active{name="example-server.rancher-dev.internal",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_host_state_updating_inactive State of defined host as reported by Rancher
# TYPE rancher_host_state_updating_inactive gauge
rancher_host_state_updating_inactive{name="example-server.rancher-dev.internal",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_service_health_status HealthState of defined service as reported by Rancher
# TYPE rancher_service_health_status gauge
rancher_service_health_status{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_service_health_status{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_service_health_status{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_service_health_status{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_service_health_status{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_service_health_status{name="prometheus-rancher-exporter",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_service_health_status{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
# HELP rancher_service_scale scale of defined service as reported by Rancher
# TYPE rancher_service_scale gauge
rancher_service_scale{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_service_scale{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_service_scale{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_service_scale{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_service_scale{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_service_scale{name="prometheus-rancher-exporter",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_service_scale{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
# HELP rancher_service_state_activating Service State of defined stack as reported by Rancher
# TYPE rancher_service_state_activating gauge
rancher_service_state_activating{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_activating{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_activating{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_activating{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_activating{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_activating{name="prometheus-rancher-exporter",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_activating{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_service_state_active Service State of defined stack as reported by Rancher
# TYPE rancher_service_state_active gauge
rancher_service_state_active{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_service_state_active{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_service_state_active{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_service_state_active{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_service_state_active{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_service_state_active{name="prometheus-rancher-exporter",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_active{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
# HELP rancher_service_state_canceled_upgrade HealthState of defined stack as reported by Rancher
# TYPE rancher_service_state_canceled_upgrade gauge
rancher_service_state_canceled_upgrade{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_canceled_upgrade{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_canceled_upgrade{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_canceled_upgrade{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_canceled_upgrade{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_canceled_upgrade{name="prometheus-rancher-exporter",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_canceled_upgrade{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_service_state_canceling_upgrade HealthState of defined stack as reported by Rancher
# TYPE rancher_service_state_canceling_upgrade gauge
rancher_service_state_canceling_upgrade{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_canceling_upgrade{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_canceling_upgrade{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_canceling_upgrade{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_canceling_upgrade{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_canceling_upgrade{name="prometheus-rancher-exporter",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_canceling_upgrade{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_service_state_deactivating HealthState of defined stack as reported by Rancher
# TYPE rancher_service_state_deactivating gauge
rancher_service_state_deactivating{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_deactivating{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_deactivating{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_deactivating{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_deactivating{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_deactivating{name="prometheus-rancher-exporter",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_deactivating{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_service_state_finishing_upgrade HealthState of defined stack as reported by Rancher
# TYPE rancher_service_state_finishing_upgrade gauge
rancher_service_state_finishing_upgrade{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_finishing_upgrade{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_finishing_upgrade{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_finishing_upgrade{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_finishing_upgrade{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_finishing_upgrade{name="prometheus-rancher-exporter",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_finishing_upgrade{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_service_state_inactive HealthState of defined stack as reported by Rancher
# TYPE rancher_service_state_inactive gauge
rancher_service_state_inactive{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_inactive{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_inactive{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_inactive{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_inactive{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_inactive{name="prometheus-rancher-exporter",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_inactive{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_service_state_registering HealthState of defined stack as reported by Rancher
# TYPE rancher_service_state_registering gauge
rancher_service_state_registering{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_registering{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_registering{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_registering{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_registering{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_registering{name="prometheus-rancher-exporter",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_registering{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_service_state_removed HealthState of defined stack as reported by Rancher
# TYPE rancher_service_state_removed gauge
rancher_service_state_removed{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_removed{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_removed{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_removed{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_removed{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_removed{name="prometheus-rancher-exporter",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_removed{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_service_state_removing HealthState of defined stack as reported by Rancher
# TYPE rancher_service_state_removing gauge
rancher_service_state_removing{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_removing{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_removing{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_removing{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_removing{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_removing{name="prometheus-rancher-exporter",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_removing{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_service_state_requested HealthState of defined stack as reported by Rancher
# TYPE rancher_service_state_requested gauge
rancher_service_state_requested{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_requested{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_requested{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_requested{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_requested{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_requested{name="prometheus-rancher-exporter",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_requested{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_service_state_restarting HealthState of defined stack as reported by Rancher
# TYPE rancher_service_state_restarting gauge
rancher_service_state_restarting{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_restarting{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_restarting{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_restarting{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_restarting{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_restarting{name="prometheus-rancher-exporter",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_restarting{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_service_state_rolling_back HealthState of defined stack as reported by Rancher
# TYPE rancher_service_state_rolling_back gauge
rancher_service_state_rolling_back{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_rolling_back{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_rolling_back{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_rolling_back{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_rolling_back{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_rolling_back{name="prometheus-rancher-exporter",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_rolling_back{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_service_state_updating_active HealthState of defined stack as reported by Rancher
# TYPE rancher_service_state_updating_active gauge
rancher_service_state_updating_active{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_updating_active{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_updating_active{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_updating_active{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_updating_active{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_updating_active{name="prometheus-rancher-exporter",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_updating_active{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_service_state_updating_inactive HealthState of defined stack as reported by Rancher
# TYPE rancher_service_state_updating_inactive gauge
rancher_service_state_updating_inactive{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_updating_inactive{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_updating_inactive{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_updating_inactive{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_updating_inactive{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_updating_inactive{name="prometheus-rancher-exporter",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_updating_inactive{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_service_state_upgraded HealthState of defined stack as reported by Rancher
# TYPE rancher_service_state_upgraded gauge
rancher_service_state_upgraded{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_upgraded{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_upgraded{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_upgraded{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_upgraded{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_upgraded{name="prometheus-rancher-exporter",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_service_state_upgraded{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_service_state_upgrading HealthState of defined stack as reported by Rancher
# TYPE rancher_service_state_upgrading gauge
rancher_service_state_upgrading{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_upgrading{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_upgrading{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_upgrading{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_upgrading{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_upgrading{name="prometheus-rancher-exporter",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_service_state_upgrading{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_stack_health_state HealthState of defined stack as reported by Rancher
# TYPE rancher_stack_health_state gauge
rancher_stack_health_state{name="Ed-Test",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_stack_health_state{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_stack_health_state{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_stack_health_state{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_stack_health_state{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_stack_health_state{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_stack_health_state{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
# HELP rancher_stack_state_activating State of defined stack as reported by Rancher
# TYPE rancher_stack_state_activating gauge
rancher_stack_state_activating{name="Ed-Test",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_activating{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_activating{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_activating{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_activating{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_activating{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_activating{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_stack_state_active State of defined stack as reported by Rancher
# TYPE rancher_stack_state_active gauge
rancher_stack_state_active{name="Ed-Test",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_stack_state_active{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_stack_state_active{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_stack_state_active{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_stack_state_active{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_stack_state_active{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
rancher_stack_state_active{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 1
# HELP rancher_stack_state_canceling_upgrade State of defined stack as reported by Rancher
# TYPE rancher_stack_state_canceling_upgrade gauge
rancher_stack_state_canceling_upgrade{name="Ed-Test",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_canceling_upgrade{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_canceling_upgrade{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_canceling_upgrade{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_canceling_upgrade{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_canceling_upgrade{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_canceling_upgrade{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_stack_state_cancelled_upgrade State of defined stack as reported by Rancher
# TYPE rancher_stack_state_cancelled_upgrade gauge
rancher_stack_state_cancelled_upgrade{name="Ed-Test",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_cancelled_upgrade{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_cancelled_upgrade{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_cancelled_upgrade{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_cancelled_upgrade{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_cancelled_upgrade{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_cancelled_upgrade{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_stack_state_error State of defined stack as reported by Rancher
# TYPE rancher_stack_state_error gauge
rancher_stack_state_error{name="Ed-Test",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_error{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_error{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_error{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_error{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_error{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_error{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_stack_state_erroring State of defined stack as reported by Rancher
# TYPE rancher_stack_state_erroring gauge
rancher_stack_state_erroring{name="Ed-Test",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_erroring{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_erroring{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_erroring{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_erroring{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_erroring{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_erroring{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_stack_state_finishing_upgrade State of defined stack as reported by Rancher
# TYPE rancher_stack_state_finishing_upgrade gauge
rancher_stack_state_finishing_upgrade{name="Ed-Test",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_finishing_upgrade{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_finishing_upgrade{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_finishing_upgrade{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_finishing_upgrade{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_finishing_upgrade{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_finishing_upgrade{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_stack_state_removed State of defined stack as reported by Rancher
# TYPE rancher_stack_state_removed gauge
rancher_stack_state_removed{name="Ed-Test",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_removed{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_removed{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_removed{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_removed{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_removed{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_removed{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_stack_state_removing State of defined stack as reported by Rancher
# TYPE rancher_stack_state_removing gauge
rancher_stack_state_removing{name="Ed-Test",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_removing{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_removing{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_removing{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_removing{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_removing{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_removing{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_stack_state_requested State of defined stack as reported by Rancher
# TYPE rancher_stack_state_requested gauge
rancher_stack_state_requested{name="Ed-Test",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_requested{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_requested{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_requested{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_requested{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_requested{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_requested{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_stack_state_rolling_back State of defined stack as reported by Rancher
# TYPE rancher_stack_state_rolling_back gauge
rancher_stack_state_rolling_back{name="Ed-Test",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_rolling_back{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_rolling_back{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_rolling_back{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_rolling_back{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_rolling_back{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_rolling_back{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_stack_state_updating_active State of defined stack as reported by Rancher
# TYPE rancher_stack_state_updating_active gauge
rancher_stack_state_updating_active{name="Ed-Test",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_updating_active{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_updating_active{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_updating_active{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_updating_active{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_updating_active{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_updating_active{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_stack_state_upgraded State of defined stack as reported by Rancher
# TYPE rancher_stack_state_upgraded gauge
rancher_stack_state_upgraded{name="Ed-Test",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_upgraded{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_upgraded{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_upgraded{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_upgraded{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_upgraded{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_upgraded{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
# HELP rancher_stack_state_upgrading State of defined stack as reported by Rancher
# TYPE rancher_stack_state_upgrading gauge
rancher_stack_state_upgrading{name="Ed-Test",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_upgrading{name="dns",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_upgrading{name="healthcheck",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_upgrading{name="ipsec",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_upgrading{name="metadata",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_upgrading{name="network-manager",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
rancher_stack_state_upgrading{name="scheduler",rancherURL="http://x.x.x.x:8080/v2-beta"} 0
```

## Metadata
[![](https://images.microbadger.com/badges/version/infinityworks/prometheus-rancher-exporter.svg)](http://microbadger.com/images/infinityworks/prometheus-rancher-exporter "Get your own version badge on microbadger.com") [![](https://images.microbadger.com/badges/image/infinityworks/prometheus-rancher-exporter.svg)](http://microbadger.com/images/infinityworks/prometheus-rancher-exporter "Get your own image badge on microbadger.com")
