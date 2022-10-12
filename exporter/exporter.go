package exporter

import (
	"github.com/fffonion/tplink-plug-exporter/kasa"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type Exporter struct {
	target string
	client *kasa.KasaClient

	metricsUp,
	metricsMetadata,
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
		labelNames  = []string{"alias", "id"}
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

		metricsMetadata: prometheus.NewDesc("kasa_metadata",
			"Device metadata.",
			[]string{
				"alias", "hw_ver", "sw_ver", "model", "feature",
			}, constLabels,
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
	s := k.client.SystemService(nil)
	r, err := s.GetSysInfo()

	if err != nil {
		ch <- prometheus.MustNewConstMetric(k.metricsUp, prometheus.GaugeValue,
			0)
		log.Errorln("error collecting", k.target, ":", err)
		return
	}

	// "alias", "hw_ver", "sw_ver", "model", "feature",
	ch <- prometheus.MustNewConstMetric(k.metricsMetadata, prometheus.GaugeValue,
		1, r.Alias, r.HardwareVersion, r.SoftwareVersion, r.Model, r.Feature)

	ch <- prometheus.MustNewConstMetric(k.metricsRelayState, prometheus.GaugeValue,
		float64(r.RelayState), r.Alias, r.DeviceID)
	ch <- prometheus.MustNewConstMetric(k.metricsOnTime, prometheus.CounterValue,
		float64(r.OnTime), r.Alias, r.DeviceID)
	ch <- prometheus.MustNewConstMetric(k.metricsRssi, prometheus.GaugeValue,
		float64(r.RSSI), r.Alias, r.DeviceID)

	aliases := map[string]string{}
	emeterContexts := []*kasa.KasaRequestContext{
		nil, // a nil context, represent the single plug or the parent strip
	}

	// iterrate over every child plug in a power strip
	for _, children := range r.Children {
		aliases[children.ID] = children.Alias
		emeterContexts = append(emeterContexts, &kasa.KasaRequestContext{
			ChildIDs: []string{children.ID},
		})

		ch <- prometheus.MustNewConstMetric(k.metricsRelayState, prometheus.GaugeValue,
			float64(children.State), children.Alias, children.ID)

		ch <- prometheus.MustNewConstMetric(k.metricsOnTime, prometheus.CounterValue,
			float64(children.OnTime), children.Alias, children.ID)
	}

	if s.EmeterSupported(r) {
		for _, ctx := range emeterContexts {
			m := k.client.EmeterService(ctx)
			re, err := m.GetRealtime()

			labels := []string{r.Alias, r.DeviceID}
			// if this is a child plug in a powerstrip, set the alias and ID
			if ctx != nil {
				id := ctx.ChildIDs[0]
				labels[0] = aliases[id]
				labels[1] = id
			}

			// TODO: only set the child up to 0 on error
			if err != nil {
				ch <- prometheus.MustNewConstMetric(k.metricsUp, prometheus.GaugeValue,
					0)
				log.Errorln("error collecting", k.target, ":", err)
				return
			}

			ch <- prometheus.MustNewConstMetric(k.metricsCurrent, prometheus.GaugeValue,
				float64(re.Current), labels...)
			ch <- prometheus.MustNewConstMetric(k.metricsVoltage, prometheus.GaugeValue,
				float64(re.Voltage), labels...)
			ch <- prometheus.MustNewConstMetric(k.metricsPowerLoad, prometheus.GaugeValue,
				float64(re.Power), labels...)
			ch <- prometheus.MustNewConstMetric(k.metricsPowerTotal, prometheus.CounterValue,
				float64(re.Total), labels...)
		}

	}

	ch <- prometheus.MustNewConstMetric(k.metricsUp, prometheus.GaugeValue,
		1)
}
