---
title: Environments
description: Environment variables for Xray Checker
---

## Subscription

### SUBSCRIPTION_URL

- CLI: `--subscription-url`
- Default: -
- Description: Subscription URL for proxy configurations

### SUBSCRIPTION_UPDATE

- CLI: `--subscription-update`
- Default: `true`
- Description: Auto-update subscription

### SUBSCRIPTION_UPDATE_INTERVAL

- CLI: `--subscription-update-interval`
- Default: `300`
- Description: Subscription update interval in seconds

## Proxy

### PROXY_CHECK_INTERVAL

- CLI: `--proxy-check-interval`
- Default: `300`
- Description: Check interval in seconds

### PROXY_CHECK_METHOD

- CLI: `--proxy-check-method`
- Default: `ip`
- Description: Check method (ip/status)

### PROXY_IP_CHECK_URL

- CLI: `--proxy-ip-check-url`
- Default: `https://api.ipify.org?format=text`
- Description: IP check service URL

### PROXY_STATUS_CHECK_URL

- CLI: `--proxy-status-check-url`
- Default: `http://cp.cloudflare.com/generate_204`
- Description: Status check URL

### PROXY_TIMEOUT

- CLI: `--proxy-timeout`
- Default: `30`
- Description: Check timeout in seconds

### SIMULATE_LATENCY

- CLI: `--simulate-latency`
- Default: `true`
- Description: Add latency to response

## Xray

### XRAY_START_PORT

- CLI: `--xray-start-port`
- Default: `10000`
- Description: Starting port for configurations

### XRAY_LOG_LEVEL

- CLI: `--xray-log-level`
- Default: `none`
- Description: Log level (debug/info/warning/error/none)

## Metrics

### METRICS_PORT

- CLI: `--metrics-port`
- Default: `2112`
- Description: Metrics port

### METRICS_PROTECTED

- CLI: `--metrics-protected`
- Default: `false`
- Description: Protect metrics with Basic Auth

### METRICS_USERNAME

- CLI: `--metrics-username`
- Default: `metricsUser`
- Description: Basic Auth username

### METRICS_PASSWORD

- CLI: `--metrics-password`
- Default: `MetricsVeryHardPassword`
- Description: Basic Auth password

### METRICS_PUSH_URL

- CLI: `--metrics-push-url`
- Default: -
- Description: Prometheus pushgateway URL

### METRICS_INSTANCE

- CLI: `--metrics-instance`
- Default: -
- Description: Instance label for metrics

## Other

### RUN_ONCE

- CLI: `--run-once`
- Default: `false`
- Description: Run one check cycle and exit
