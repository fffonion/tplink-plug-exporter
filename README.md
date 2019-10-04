# tplink-plug-exporter

## Sample prometheus config

```yaml
# scrape kasa devices
scrape_configs:
  - job_name: 'kasa'
    static_configs:
    - targets:
      - 192.168.0.233
      - 192.168.0.234
    metrics_path: /scrape
    relabel_configs:
      - source_labels : [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: localhost:9233

# scrape kasa_exporter itself
  - job_name: 'kasa_exporter'
    static_configs:
      - targets:
        - localhost:9233
```