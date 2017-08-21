FROM golang:1.8-alpine as builder
MAINTAINER infinityworks

EXPOSE 9173

COPY . /go/src/github.com/infinityworks/prometheus-rancher-exporter

RUN apk --update add ca-certificates \
 && apk --update add --virtual build-deps go git \
 && cd /go/src/github.com/infinityworks/prometheus-rancher-exporter \
 && GOPATH=/go go get \
 && GOPATH=/go go build -o /bin/rancher_exporter \
 && apk del --purge build-deps \
 && rm -rf /go/bin /go/pkg /var/cache/apk/*


FROM alpine

ENV LISTEN_PORT=9173

RUN addgroup exporter \
     && adduser -S -G exporter exporter \
     && apk --update --no-cache add ca-certificates

COPY --from=builder /bin/rancher_exporter /bin/rancher_exporter

ENTRYPOINT [ "/bin/rancher_exporter" ]
