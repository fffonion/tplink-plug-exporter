package main

import (
	"net/http"

	"github.com/fffonion/tplink-plug-exporter/exporter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	s := exporter.NewHttpServer()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9233", s)
}
