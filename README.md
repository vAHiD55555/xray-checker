# Xray Checker

[![GitHub Release](https://img.shields.io/github/v/release/kutovoys/xray-checker?style=flat&color=blue)](https://github.com/kutovoys/xray-checker/releases/latest)
[![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/kutovoys/xray-checker/build-publish.yml)](https://github.com/kutovoys/xray-checker/actions/workflows/build-publish.yml)
[![DockerHub](https://img.shields.io/badge/DockerHub-kutovoys%2Fxray--checker-blue)](https://hub.docker.com/r/kutovoys/xray-checker/)
[![GitHub License](https://img.shields.io/github/license/kutovoys/xray-checker?color=greeen)](https://github.com/kutovoys/xray-checker/blob/main/LICENSE)
[![en](https://img.shields.io/badge/lang-en-red)](https://github.com/kutovoys/xray-checker/blob/main/README.md)
[![ru](https://img.shields.io/badge/lang-ru-blue)](https://github.com/kutovoys/xray-checker/blob/main/README_RU.md)

Xray Checker is a tool for monitoring proxy server availability, supporting VLESS, Trojan, and Shadowsocks protocols. It automatically tests connections through Xray Core and provides metrics for Prometheus, as well as API endpoints for integration with monitoring systems.

<div align="center">
  <img src="images/xray-checker.png" alt="Dashboard Screenshot">
</div>

## Features

- **Protocol Support**: VLESS, Trojan, and Shadowsocks
- **Prometheus Metrics**: Export proxy status metrics for Prometheus
- **API Endpoints**: Individual endpoints for each proxy for monitoring system integration
- **Automatic Updates**: Periodic checking and updating of configuration from subscription URL
- **Web Interface**: Simple interface for viewing proxy status and configuration
- **Basic Auth**: Optional protection of metrics and API using basic authentication
- **Docker Support**: Easy deployment using Docker and Docker Compose

## Metrics

The exporter provides the following metrics:

| Name                | Description                               |
| ------------------- | ----------------------------------------- |
| `xray_proxy_status` | Proxy status (1: working, 0: not working) |

Each metric includes the following labels:

- `protocol`: Protocol type (vless/trojan/shadowsocks)
- `address`: Server address and port
- `name`: Proxy configuration name

## Configuration

| Environment Variable           | CLI Argument                     | Default                                 | Description                               |
| ------------------------------ | -------------------------------- | --------------------------------------- | ----------------------------------------- |
| `SUBSCRIPTION_URL`             | `--subscription-url`             | -                                       | Subscription URL for proxy configurations |
| `SUBSCRIPTION_UPDATE`          | `--subscription-update`          | `true`                                  | Auto-update subscription                  |
| `SUBSCRIPTION_UPDATE_INTERVAL` | `--subscription-update-interval` | `300`                                   | Subscription update interval in seconds   |
|                                |                                  |                                         |                                           |
| `PROXY_CHECK_INTERVAL`         | `--proxy-check-interval`         | `300`                                   | Check interval in seconds                 |
| `PROXY_CHECK_METHOD`           | `--proxy-check-method`           | `ip`                                    | Check method (ip/status)                  |
| `PROXY_IP_CHECK_URL`           | `--proxy-ip-check-url`           | `https://api.ipify.org?format=text`     | IP check service URL                      |
| `PROXY_STATUS_CHECK_URL`       | `--proxy-status-check-url`       | `http://cp.cloudflare.com/generate_204` | Status check URL                          |
| `PROXY_TIMEOUT`                | `--proxy-timeout`                | `30`                                    | Check timeout in seconds                  |
| `SIMULATE_LATENCY`             | `--simulate-latency`             | `true`                                  | Add latency to response                   |
|                                |                                  |                                         |                                           |
| `XRAY_START_PORT`              | `--xray-start-port`              | `10000`                                 | Starting port for configurations          |
| `XRAY_LOG_LEVEL`               | `--xray-log-level`               | `none`                                  | Log level (debug/info/warning/error/none) |
|                                |                                  |                                         |                                           |
| `METRICS_PORT`                 | `--metrics-port`                 | `2112`                                  | Metrics port                              |
| `METRICS_PROTECTED`            | `--metrics-protected`            | `false`                                 | Protect metrics with Basic Auth           |
| `METRICS_USERNAME`             | `--metrics-username`             | `metricsUser`                           | Basic Auth username                       |
| `METRICS_PASSWORD`             | `--metrics-password`             | `MetricsVeryHardPassword`               | Basic Auth password                       |

### Subscription Format

The content of `SUBSCRIPTION_URL` must be in Base64 Encoded format containing a list of proxies. (Standard format for Xray clients - Streisand, V2rayNG).

Proxies with ports 0, 1 will be ignored.

Request headers sent:

```
Accept: */*
User-Agent: Xray-Checker
```

## Usage

### CLI

```bash
# Basic usage
./xray-checker --subscription-url="https://your-subscription-url/sub"
```

```bash
# Advanced usage with custom settings
./xray-checker \
  --subscription-url="https://your-subscription-url/sub" \
  --proxy-check-interval=300 \
  --proxy-timeout=5 \
  --proxy-check-method=ip \
  --proxy-ip-check-url="https://api.ipify.org?format=text" \
  --proxy-status-check-url="http://cp.cloudflare.com/generate_204" \
  --simulate-latency=true \
  --xray-start-port=10000 \
  --xray-log-level=none \
  --metrics-port=2112 \
  --metrics-protected=true \
  --metrics-username=custom_user \
  --metrics-password=custom_pass
```

### Docker

```bash
docker run -d \
  -e SUBSCRIPTION_URL=https://your-subscription-url/sub \
  -e PROXY_CHECK_INTERVAL=300 \
  -e PROXY_CHECK_METHOD=ip \
  -e PROXY_TIMEOUT=30 \
  -e XRAY_START_PORT=10000 \
  -e METRICS_PORT=2112 \
  kutovoys/xray-checker
```

### Docker Compose

```yaml
services:
  xray-checker:
    image: kutovoys/xray-checker
    environment:
      - SUBSCRIPTION_URL=https://your-subscription-url/sub
      - PROXY_CHECK_INTERVAL=300
      - PROXY_CHECK_METHOD=ip
      - PROXY_TIMEOUT=30
      - XRAY_START_PORT=10000
      - METRICS_PROTECTED=true
      - METRICS_USERNAME=custom_user
      - METRICS_PASSWORD=custom_password
    ports:
      - "2112:2112"
```

### Prometheus Configuration

Add the following to your prometheus.yml:

```yaml
scrape_configs:
  - job_name: "xray-checker"
    static_configs:
      - targets: ["localhost:2112"]
    scrape_interval: 1m
```

## API Endpoints

- `/` - Information page
- `/metrics` - Prometheus metrics endpoint
- `/health` - Health check endpoint
- `/config/{protocol}-{address}-{port}` - Status of specific proxy (returns 200 OK if working, 503 if not)

### Uptime Kuma Integration

You can monitor each proxy using its dedicated endpoint in Uptime Kuma:

1. Add new monitor
2. Select "HTTP(s)"
3. Enter URL: `http://your-server:2112/config/vless-example.com-443`
4. The monitor will show "Up" when proxy is working and "Down" when it fails

## Connection Check Logic

1. Initial setup:

   - Retrieve configurations from subscription URL
   - Generate unified Xray configuration file
   - Start Xray Core instance

2. Periodic checks (every N minutes):
   - Get current IP without proxy
   - For each proxy configuration:
     - Connect through local SOCKS5 port
     - Try to get IP through proxy
     - Compare IPs to determine if proxy is working
     - Update Prometheus metrics and internal status
   - Check subscription URL for changes
     - If changes detected:
       - Generate new Xray configuration
       - Restart Xray Core instance
       - Update endpoints

## Contribute

Contributions to Xray Checker are warmly welcomed. Whether it's bug fixes, new features, or documentation improvements, here's a quick guide to contributing:

1. **Fork & Branch**: Fork this repository and create a branch for your changes
2. **Implement**: Make your changes while keeping code clean and documented
3. **Test**: Ensure your changes don't break existing functionality
4. **Commit & PR**: Create commits with clear messages and open a Pull Request
5. **Feedback**: Be ready to engage with feedback and refine your contribution

If you're new to this, GitHub provides an excellent guide on [creating a pull request](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request).
