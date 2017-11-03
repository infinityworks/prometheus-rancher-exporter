FROM golang:1.9.1-alpine3.6
LABEL maintainer="Infinity Works"

EXPOSE 9173

RUN addgroup exporter \
 && adduser -S -G exporter exporter

COPY . /go/src/github.com/infinityworks/prometheus-rancher-exporter

RUN apk --update add ca-certificates \
 && apk --update add --virtual build-deps go git \
 && cd /go/src/github.com/infinityworks/prometheus-rancher-exporter \
 && GOPATH=/go go get \
 && GOPATH=/go go build -o /bin/rancher_exporter \
 && apk del --purge build-deps \
 && rm -rf /go/bin /go/pkg /var/cache/apk/*

USER exporter

ENTRYPOINT [ "/bin/rancher_exporter" ]
