FROM golang:1.12-alpine AS builder
ADD . /src
RUN apk add --no-cache git
WORKDIR /src
RUN go build main.go

FROM alpine:latest
COPY --from=builder /src/main /tplink-plug-exporter
EXPOSE 9233
ENTRYPOINT ["/tplink-plug-exporter"]
