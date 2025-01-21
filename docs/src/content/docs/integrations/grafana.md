---
title: Grafana Dashboards
description: Grafana integration
---

### Basic Dashboard

Import this dashboard for basic monitoring:

```json
{
  "title": "Xray Checker Overview",
  "tags": ["xray-checker"],
  "timezone": "browser",
  "panels": [
    {
      "title": "Proxy Status",
      "type": "stat",
      "targets": [
        {
          "expr": "xray_proxy_status",
          "legendFormat": "{{name}}"
        }
      ]
    },
    {
      "title": "Proxy Latency",
      "type": "graph",
      "targets": [
        {
          "expr": "xray_proxy_latency_ms",
          "legendFormat": "{{name}}"
        }
      ]
    }
  ]
}
```

### Advanced Dashboard

Features:

- Status heatmap
- Latency trends
- Instance comparison

Available at:
