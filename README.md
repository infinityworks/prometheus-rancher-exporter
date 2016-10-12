# prometheus-rancher-exporter

Exposes the health of stacks/services and hosts from the Rancher API, to a Prometheus compatible endpoint. 

## Description

This container makes use of Ranchers ability to assign API access to a container at runtime. This is achieved through labels to create a connection to the API.
The application, expects to get the following environment variables from the host, if not using the supplied labelss in rancher-compose then you can update these values yourself, using environment variables.

* CATTLE_ACCESS_KEY
* CATTLE_SECRET_KEY
* CATTLE_URL

## Install and deploy

Run manually from Docker Hub:
```
docker run -d --restart=always -p 9010:9010 infinityworks/prometheus-rancher-exporter
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
    stdin_open: true
    labels:
      io.rancher.container.create_agent: true
      io.rancher.container.agent.role: environment
    expose:
      - 9010:9010
    image: infinityworks/prometheus-rancher-exporter:latest
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

## Metadata
[![](https://images.microbadger.com/badges/version/infinityworks/prometheus-rancher-exporter.svg)](http://microbadger.com/images/infinityworks/prometheus-rancher-exporter "Get your own version badge on microbadger.com") [![](https://images.microbadger.com/badges/image/infinityworks/prometheus-rancher-exporter.svg)](http://microbadger.com/images/infinityworks/prometheus-rancher-exporter "Get your own image badge on microbadger.com")
