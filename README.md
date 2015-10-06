prometheus-rancher-exporter
===========================

Exposes Rancher environment status to Prometheus.

## Install and deploy

Run from Docker Hub:
```
docker run -d --restart=always -p 9010:9010 -e HOST=<host> -e PORT=8080 -e API_ACCESS_KEY=<api-key> -e API_SECRET_KEY=<secret-key> barwell/prometheus-rancher-exporter
```

Build a docker image:
```
docker build -t <image-name> .
docker run -d --restart=always -p 9010:9010 -e HOST=<host> -e PORT=8080 -e API_ACCESS_KEY=<api-key> -e API_SECRET_KEY=<secret-key> <image-name>
```

Running the node process:
```
HOST=<host> PORT=8080 API_ACCESS_KEY=<api-key> API_SECRET_KEY=<secret-key> DEBUG=re node app.js
```

## Metrics

Metrics will be made available on port 9010 by default, or you can pass environment variable ```LISTEN_PORT``` to override this.

```
# HELP rancher_environment Value of 1 if all containers in a stack are active
# TYPE rancher_environment gauge
rancher_environment{name="test1"} 1
rancher_environment{name="test2"} 0
rancher_environment{name="load_test"} 1
rancher_environment{name="preprod"} 1
```
