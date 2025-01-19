---
title: Prometheus
description: Prometheus integration options and examples
---

Xray Checker supports two ways to integrate with Prometheus:

## Direct metrics collection

Add the following configuration to your prometheus.yml:

```yaml
scrape_configs:
  - job_name: "xray-checker"
    static_configs:
      - targets: ["xray-checker:8080"]
```

## Pushgateway integration

Pushgateway integration is useful when you need to collect metrics from multiple instances or don't have access to the server where Xray Checker is running.
