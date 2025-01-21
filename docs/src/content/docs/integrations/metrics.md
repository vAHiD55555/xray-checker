---
title: Metrics
description: Metrics options and examples
---

Xray Checker provides two Prometheus metrics for monitoring proxy status and performance:

### xray_proxy_status

Status metric indicating proxy availability:

- Type: Gauge
- Values: 1 (working) or 0 (failed)
- Labels:
  - `protocol`: Proxy protocol (vless/vmess/trojan/shadowsocks)
  - `address`: Server address and port
  - `name`: Proxy configuration name
  - `instance`: Instance name (if configured)

Example:

```text
# HELP xray_proxy_status Status of proxy connection (1: success, 0: failure)
# TYPE xray_proxy_status gauge
xray_proxy_status{protocol="vless",address="example.com:443",name="proxy1",instance="dc1"} 1
```

### xray_proxy_latency_ms

Latency metric showing connection response time:

- Type: Gauge
- Values: Milliseconds (0 if failed)
- Labels: Same as xray_proxy_status

Example:

```text
# HELP xray_proxy_latency_ms Latency of proxy connection in milliseconds
# TYPE xray_proxy_latency_ms gauge
xray_proxy_latency_ms{protocol="vless",address="example.com:443",name="proxy1",instance="dc1"} 156
```
