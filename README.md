prometheus-rancher-exporter
===========================

Exposes Rancher metrics to Prometheus

## Install and deploy

Run from Docker Hub:
```
docker run -d --restart=always -e HOST=<host> -e PORT=8080 -e API_ACCESS_KEY=<api-key> -e API_SECRET_KEY=<secret-key> barwell/prometheus-rancher-exporter
```

Build a docker image:
```
docker build -t <image-name> .
docker run -d --restart=always -e HOST=<host> -e PORT=8080 -e API_ACCESS_KEY=<api-key> -e API_SECRET_KEY=<secret-key> <image-name>
```

With npm:
```
npm install prometheus-rancher-exporter
HOST=<host> PORT=8080 API_ACCESS_KEY=<api-key> API_SECRET_KEY=<secret-key> DEBUG=re node app.js
```

## Metrics

Metrics will be made available on port 9010 by default, or you can pass environment variable ```LISTEN_PORT``` to override this.
