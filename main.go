package main

import (
	"flag"
	"net/http"

	"github.com/prometheus/common/log"

	"github.com/fffonion/tplink-plug-exporter/exporter"
)

func main() {
	var metricsAddr = flag.String("metrics.listen-addr", ":9233", "listen address for tplink-plug exporter")

	flag.Parse()
	s := exporter.NewHttpServer()
	log.Infof("Accepting Prometheus Requests on %s", *metricsAddr)
	log.Fatal(http.ListenAndServe(*metricsAddr, s))
}
