---
title: Grafana
description: Примеры панелей Grafana
---

## Панели Grafana

### Статус прокси

Панель статуса с использованием метрики `xray_proxy_status`:

```jsonc
{
  "targets": [
    {
      "expr": "xray_proxy_status",
      "legendFormat": "{{name}}"
    }
  ],
  "type": "stat",
  "options": {
    "colorMode": "value",
    "graphMode": "none",
    "justifyMode": "auto",
    "textMode": "auto",
    "colorThresholds": [
      { "value": 0, "color": "red" },
      { "value": 1, "color": "green" }
    ]
  }
}
```

### График задержки

График задержки с использованием метрики `xray_proxy_latency_ms`:

```jsonc
{
  "targets": [
    {
      "expr": "xray_proxy_latency_ms",
      "legendFormat": "{{name}}"
    }
  ],
  "type": "timeseries",
  "options": {
    "tooltip": {
      "mode": "multi",
      "sort": "desc"
    }
  },
  "fieldConfig": {
    "defaults": {
      "color": {
        "mode": "palette-classic"
      },
      "custom": {
        "fillOpacity": 10,
        "lineWidth": 1,
        "spanNulls": false
      },
      "unit": "ms"
    }
  }
}
```

### Таблица прокси

Таблица со статусом и задержкой:

```jsonc
{
  "targets": [
    {
      "expr": "xray_proxy_status",
      "format": "table",
      "instant": true
    },
    {
      "expr": "xray_proxy_latency_ms",
      "format": "table",
      "instant": true
    }
  ],
  "type": "table",
  "options": {
    "showHeader": true,
    "footer": {
      "show": false
    }
  },
  "fieldConfig": {
    "defaults": {
      "custom": {
        "align": "left",
        "width": 100
      }
    },
    "overrides": [
      {
        "matcher": {
          "id": "byName",
          "options": "Value"
        },
        "properties": [
          {
            "id": "custom.width",
            "value": 150
          },
          {
            "id": "unit",
            "value": "ms"
          }
        ]
      }
    ]
  }
}
```

### Оповещения

Правило оповещения о недоступности прокси:

```yaml
groups:
  - name: xray-checker
    rules:
      - alert: ProxyDown
        expr: xray_proxy_status == 0
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Прокси {{ $labels.name }} недоступен"
          description: "Прокси {{ $labels.name }} ({{ $labels.address }}) недоступен более 5 минут"
```

Правило оповещения о высокой задержке:

```yaml
groups:
  - name: xray-checker
    rules:
      - alert: HighLatency
        expr: xray_proxy_latency_ms > 1000
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Высокая задержка прокси {{ $labels.name }}"
          description: "Задержка прокси {{ $labels.name }} превышает 1000мс более 5 минут"
```
