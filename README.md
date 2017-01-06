# prometheus-rancher-exporter

Exposes the health of Stacks / Services and Hosts from the Rancher API, to a Prometheus compatible endpoint.


## Description

The application can be run in a number of ways, the main consumption is the Docker hub image `infinityworksltd/prometheus-rancher-exporter`. 

The application requires at a minimum, the URL of the Rancher API. If you have authentication enabled on your Rancher server, the application will require a `RANCHER_ACCESS_KEY` and a `RANCHER_SECRET_KEY` providing.

If you are running the application in a Rancher managed container, you can make use of Rancher labels  to obtain an API key and auto-provision all of this information, details of this can be seen in the Docker Compose section.

If you are using this externally to Rancher, or without the use of the labels to obtain an API key, you can update these values yourself, using environment variables.

**Required**
* `CATTLE_URL` // Either provisioned through labels, or set by the user. Should be in a format similar to `http://<YOUR_IP>:8080/v2-beta`.

**Optional**
* `CATTLE_ACCESS_KEY`   // Rancher API access Key, if supplied this will be used when authentication is enabled.
* `CATTLE_SECRET_KEY`   // Rancher API secret Key, if supplied this will be used when authentication is enabled.
* `METRICS_PATH`        // Path under which to expose metrics.
* `LISTEN_ADDRESS`      // Port on which to expose metrics.
* `HIDE_SYS`            // If set to `true` then this hides any of Ranchers internal system services from being shown.
*	`LOG_LEVEL`           // Optional - Set the logging level, defaults to Info

## Compatibility

Along with the release of Rancher 1.2, a new API was introduced, the oppertunity was taken to re-write the exporter into Golang, so it's more comparible to the platforms it's interacting with. 
Testing has focused on the `v1` and `v2-beta` available with Rancher 1.2.  The `v1` support should in theory work on older versions of Rancher Server but testing has been limited.

If you find any issues, bug reports or PR's are more than welcome.

## Install and deploy

Run manually from Docker Hub:
```
docker run -d -e CATTLE_ACCESS_KEY="XXXXXXXX" -e CATTLE_SECRET_KEY="XXXXXXX" -e CATTLE_URL="http://<YOUR_IP>:8080/v2-beta" -p 9173:9173 infinityworks/prometheus-rancher-exporter
```

Build a docker image:
```
docker build -t <image-name> .
docker run -d -e CATTLE_ACCESS_KEY="XXXXXXXX" -e CATTLE_SECRET_KEY="XXXXXXX" -e CATTLE_URL="http://<YOUR_IP>:8080/v2-beta" -p 9173:9173 <image-name>
```

## Docker compose

For users running the container within a Rancher managed environment:
```
prometheus-rancher-exporter:
    tty: true
    stdin_open: true
    labels:
      io.rancher.container.create_agent: true
      io.rancher.container.agent.role: environment
    expose:
      - 9173:9173
    image: infinityworks/prometheus-rancher-exporter:latest
```

For users running the container outside a Rancher managed environment:
```
prometheus-rancher-exporter:
    tty: true
    stdin_open: true
    environment:
      - CATTLE_ACCESS_KEY="xxxx"
      - CATTLE_SECRET_KEY="xxxxxx"
      - CATTLE_URL="http://<YOUR_IP>:8080/v2-beta"
      - HIDE_SYS=true
    expose:
      - 9173:9173
    image: infinityworks/prometheus-rancher-exporter:latest
```


## Metrics

Metrics will be made available on port 9173 by default, or you can pass environment variable ```LISTEN_ADDRESS``` to override this.
An example printout of the metrics you should expect to see can be found in `METRICS.md`.


## Metadata
[![](https://images.microbadger.com/badges/version/infinityworks/prometheus-rancher-exporter.svg)](http://microbadger.com/images/infinityworks/prometheus-rancher-exporter "Get your own version badge on microbadger.com") [![](https://images.microbadger.com/badges/image/infinityworks/prometheus-rancher-exporter.svg)](http://microbadger.com/images/infinityworks/prometheus-rancher-exporter "Get your own image badge on microbadger.com")
[![Go Report Card](https://goreportcard.com/badge/github.com/infinityworksltd/prometheus-rancher-exporter)](https://goreportcard.com/report/github.com/infinityworksltd/prometheus-rancher-exporter)
[![GoDoc](https://godoc.org/github.com/infinityworksltd/prometheus-rancher-exporter?status.svg)](https://godoc.org/github.com/infinityworksltd/prometheus-rancher-exporter)