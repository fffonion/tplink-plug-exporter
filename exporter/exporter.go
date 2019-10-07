package exporter

import (
	"github.com/fffonion/tplink-plug-exporter/kasa"
	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	target string
	client *kasa.KasaClient

	metricsUp,
	metricsRelayState,
	metricsOnTime,
	metricsRssi,
	metricsCurrent,
	metricsVoltage,
	metricsPowerLoad,
	metricsPowerTotal *prometheus.Desc
}

type ExporterTarget struct {
	Host string
}

func NewExporter(t *ExporterTarget) *Exporter {
	var (
		constLabels = prometheus.Labels{}
		labelNames  = []string{"alias"}
	)

	e := &Exporter{
		target: t.Host,
		client: kasa.New(&kasa.KasaClientConfig{
			Host: t.Host,
		}),
		metricsUp: prometheus.NewDesc("kasa_online",
			"Device online.",
			nil, constLabels,
		),
		metricsRelayState: prometheus.NewDesc("kasa_relay_state",
			"Relay state (switch on/off).",
			labelNames, constLabels,
		),
		metricsOnTime: prometheus.NewDesc("kasa_on_time",
			"Time in seconds since online.",
			labelNames, constLabels),
		metricsRssi: prometheus.NewDesc("kasa_rssi",
			"Wifi received signal strength indicator.",
			labelNames, constLabels),
		metricsCurrent: prometheus.NewDesc("kasa_current",
			"Current flowing through device in Ampere.",
			labelNames, constLabels),
		metricsVoltage: prometheus.NewDesc("kasa_voltage",
			"Current voltage connected to device in Volt.",
			labelNames, constLabels),
		metricsPowerLoad: prometheus.NewDesc("kasa_power_load",
			"Current power in Watt.",
			labelNames, constLabels),
		metricsPowerTotal: prometheus.NewDesc("kasa_power_total",
			"Power consumption since device connected in kWh.",
			labelNames, constLabels),
	}
	return e
}

func (k *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- k.metricsUp
	ch <- k.metricsRelayState
	ch <- k.metricsOnTime
	ch <- k.metricsRssi
	ch <- k.metricsCurrent
	ch <- k.metricsVoltage
	ch <- k.metricsPowerLoad
	ch <- k.metricsPowerTotal

}

func (k *Exporter) Collect(ch chan<- prometheus.Metric) {
	s := k.client.SystemService()
	r, err := s.GetSysInfo()

	if err != nil {
		ch <- prometheus.MustNewConstMetric(k.metricsUp, prometheus.GaugeValue,
			0)
		return
	}
	alias := r.Alias

	ch <- prometheus.MustNewConstMetric(k.metricsRelayState, prometheus.GaugeValue,
		float64(r.RelayState), alias)
	ch <- prometheus.MustNewConstMetric(k.metricsOnTime, prometheus.CounterValue,
		float64(r.OnTime), alias)
	ch <- prometheus.MustNewConstMetric(k.metricsRssi, prometheus.GaugeValue,
		float64(r.RSSI), alias)

	if s.EmeterSupported(r) {
		m := k.client.EmeterService()
		r, err := m.GetRealtime()

		if err != nil {
			ch <- prometheus.MustNewConstMetric(k.metricsUp, prometheus.GaugeValue,
				0)
			return
		}
		ch <- prometheus.MustNewConstMetric(k.metricsCurrent, prometheus.GaugeValue,
			float64(r.Current), alias)
		ch <- prometheus.MustNewConstMetric(k.metricsVoltage, prometheus.GaugeValue,
			float64(r.Voltage), alias)
		ch <- prometheus.MustNewConstMetric(k.metricsPowerLoad, prometheus.GaugeValue,
			float64(r.Power), alias)
		ch <- prometheus.MustNewConstMetric(k.metricsPowerTotal, prometheus.CounterValue,
			float64(r.Total), alias)
	}

	ch <- prometheus.MustNewConstMetric(k.metricsUp, prometheus.GaugeValue,
		1)
}
