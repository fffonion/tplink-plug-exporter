package main

import (
	"net/http"

	"github.com/prometheus/common/log"

	"github.com/fffonion/tplink-plug-exporter/exporter"
)

func main() {
	s := exporter.NewHttpServer()
	log.Infoln("Accepting Prometheus Requests on :9233")
	log.Fatal(http.ListenAndServe(":9233", s))
}
