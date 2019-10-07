# tplink-plug-exporter

Export TP-Link Smart Plug metrics to grafana dashboard

## Grafana dashboard

Search for `Kasa` inside grafana or install from https://grafana.com/grafana/dashboards/10957
![img](https://grafana.com/api/dashboards/10957/images/6954/image)

## Sample prometheus config

```yaml
# scrape kasa devices
scrape_configs:
  - job_name: 'kasa'
    static_configs:
    - targets:
      # IP of your smart plugs
      - 192.168.0.233
      - 192.168.0.234
    metrics_path: /scrape
    relabel_configs:
      - source_labels : [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        # IP of the exporter
        replacement: localhost:9233

# scrape kasa_exporter itself
  - job_name: 'kasa_exporter'
    static_configs:
      - targets:
        # IP of the exporter
        - localhost:9233
```

## See also

- Original reverse engineering work: https://github.com/softScheck/tplink-smartplug
