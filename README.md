prometheus-rancher-exporter
===========================

Exposes Rancher environment status to Prometheus. Makes use of Rancher labels to create a connection to the API.
Expects to get the following environment variables from the host, if not using rancher-compose then you can update these yourself:

* CATTLE_ACCESS_KEY
* CATTLE_SECRET_KEY
* CATTLE_CONFIG_URL

## Install and deploy

Run from Docker Hub:
```
docker run -d --restart=always -p 9010:9010 jolyonbrown/prometheus-rancher-exporter
```

Build a docker image:
```
docker build -t <image-name> .
docker run -d --restart=always -p 9010:9010 <image-name>
```

Running the node process:
```
DEBUG=re node app.js
```

## Docker compose

```
prometheus-rancher-exporter:
    tty: true
    labels:
      io.rancher.container.create_agent: true
      io.rancher.container.agent.role: environment
    expose:
      - 9010:9010
    image: jolyonbrown/prometheus-rancher-exporter
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
