# prometheus-rancher-exporter

Exposes the health of stacks/services and hosts from the Rancher API, to a Prometheus compatible endpoint.


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

## Compatibility

Version 1.2 of Rancher introduced a new API, we took the oppertunity to re-write the exporter into Golang so it's more comparible to the platforms it's interacting with. 
We've tested the exporter against both the V1 and V2-Beta API's available in Rancher 1.2, it should in theory work on older versions but we haven't had the chance to test. If you find any issues, bug reports or PR's are more than welcome.

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
An example printout of the metrics you should expect to see can be found in Metrics.md.


## Metadata
[![](https://images.microbadger.com/badges/version/infinityworks/prometheus-rancher-exporter.svg)](http://microbadger.com/images/infinityworks/prometheus-rancher-exporter "Get your own version badge on microbadger.com") [![](https://images.microbadger.com/badges/image/infinityworks/prometheus-rancher-exporter.svg)](http://microbadger.com/images/infinityworks/prometheus-rancher-exporter "Get your own image badge on microbadger.com")
