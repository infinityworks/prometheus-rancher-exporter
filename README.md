## Metadata
[![](https://images.microbadger.com/badges/version/infinityworks/prometheus-rancher-exporter.svg)](http://microbadger.com/images/infinityworks/prometheus-rancher-exporter "Get your own version badge on microbadger.com") [![](https://images.microbadger.com/badges/image/infinityworks/prometheus-rancher-exporter.svg)](http://microbadger.com/images/infinityworks/prometheus-rancher-exporter "Get your own image badge on microbadger.com")
[![Go Report Card](https://goreportcard.com/badge/github.com/infinityworks/prometheus-rancher-exporter)](https://goreportcard.com/report/github.com/infinityworks/prometheus-rancher-exporter)
[![GoDoc](https://godoc.org/github.com/infinityworks/prometheus-rancher-exporter?status.svg)](https://godoc.org/github.com/infinityworks/prometheus-rancher-exporter)


# prometheus-rancher-exporter

This metric exporter exposes the health of Services, Stacks and Hosts from your Rancher installation to a Prometheus endpoint.
It's designed to be run from within a Rancher managed container, the exporter translates the data available through the rancher-metadata service.

## Description

The application can be run in a number of ways, though primarily this is through usage of the Docker hub image `infinityworksltd/prometheus-rancher-exporter`. 
Configured to work 'out of the box', the exporter can also be configured through environment variables listed below.

If you are using this externally to Rancher, you would need to provide access to the rancher-metadata service and update the `METADATA_URL` value.


**Optional**
* `METADATA_URL`          // The URL of the rancher-metadata service, in a default installation this is available to the container on `http://rancher-metadata/`
* `METRICS_PATH`          // Path under which to expose metrics. Defaults to `/metrics`
* `LISTEN_PORT    `       // Port on which to expose metrics. Defaults to `9173`
* `HIDE_SYS`              // If set to `true` then this hides any of Ranchers internal system services from being shown. *If used, ensure `false` is encapsulated with quotes e.g. `HIDE_SYS="false"`.
*	`LOG_LEVEL`           // Optional - Set the logging level, defaults to `Info`

## Compatibility

`v0.26.04` and above is a substantial update and is designed to work against the rancher-metadata service. This requires changes that have been made to the latest Rancher server version.

`v0.22.52` is the last stable version that works against the API. This is designed to work with versions 1 and 2 of the rancher API found in Rancher server 1.2 +

## Installation and deployment

Run manually from Docker Hub:
```
docker run -d -p 9173:9173 infinityworks/prometheus-rancher-exporter:latest
```

Build a docker image:
```
docker build -p 9173:9173 <image-name>
```

## Docker compose

For users running the container within a Rancher managed environment:
```
prometheus-rancher-exporter:
    tty: true
    stdin_open: true
    expose:
      - 9173:9173
    image: infinityworks/prometheus-rancher-exporter:latest
```

## Metrics

Metrics will be made available on port 9173 by default, or you can pass environment variable ```LISTEN_PORT``` to override this.
An example printout of the metrics you should expect to see can be found in `METRICS.md`.

## Contributing

If you find any issues, bug reports or PR's are more than welcome. Read more in the CONTRIBUTERS document.
