---
title: Metrics
description: Metrics options and examples
---

Xray Checker provides Prometheus-format metrics for monitoring proxy server status.

## Available Metrics

### xray_proxy_status

Proxy status (1: working, 0: not working)

### xray_proxy_latency_ms

Proxy latency in milliseconds

## Metric Labels

Each metric includes the following labels:

- `protocol`: Protocol type (vless/trojan/shadowsocks)
- `address`: Server address and port
- `name`: Proxy configuration name
- `instance`: Optional instance label (if specified via --metrics-instance)
