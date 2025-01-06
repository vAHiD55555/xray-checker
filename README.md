# Xray Checker

[![GitHub Release](https://img.shields.io/github/v/release/kutovoys/xray-checker?style=flat&color=blue)](https://github.com/kutovoys/xray-checker/releases/latest)
[![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/kutovoys/xray-checker/build-publish.yml)](https://github.com/kutovoys/xray-checker/actions/workflows/build-publish.yml)
[![DockerHub](https://img.shields.io/badge/DockerHub-kutovoys%2Fxray--checker-blue)](https://hub.docker.com/r/kutovoys/xray-checker/)
[![GitHub License](https://img.shields.io/github/license/kutovoys/xray-checker?color=greeen)](https://github.com/kutovoys/xray-checker/blob/main/LICENSE)

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

The application can be configured using environment variables or command-line arguments:

| Environment Variable   | Command-Line Argument    | Required | Default                             | Description                                        |
| ---------------------- | ------------------------ | -------- | ----------------------------------- | -------------------------------------------------- |
| `SUBSCRIPTION_URL`     | `--subscription-url`     | Yes      | -                                   | Subscription URL for obtaining configurations      |
| `RECHECK_SUBSCRIPTION` | `--recheck-subscription` | No       | `true`                              | Recheck subscription on each check                 |
| `CHECK_INTERVAL`       | `--check-interval`       | No       | `300`                               | Check interval in seconds                          |
| `IP_CHECK_SERVICE`     | `--ip-check-service`     | No       | `https://api.ipify.org?format=text` | Service for IP checking                            |
| `IP_CHECK_TIMEOUT`     | `--ip-check-timeout`     | No       | `5`                                 | Timeout for IP checking in seconds                 |
| `START_PORT`           | `--start-port`           | No       | `10000`                             | Starting port for proxy configurations             |
| `XRAY_LOG_LEVEL`       | `--xray-log-level`       | No       | `none`                              | Xray logging level (debug/info/warning/error/none) |
| `METRICS_PORT`         | `--metrics-port`         | No       | `2112`                              | Port for metrics                                   |
| `METRICS_PROTECTED`    | `--metrics-protected`    | No       | `false`                             | Protect metrics with Basic Auth                    |
| `METRICS_USERNAME`     | `--metrics-username`     | No       | `metricsUser`                       | Username for Basic Auth                            |
| `METRICS_PASSWORD`     | `--metrics-password`     | No       | `MetricsVeryHardPassword`           | Password for Basic Auth                            |

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
  --check-interval=300 \
  --ip-check-timeout=5 \
  --metrics-port=2112 \
  --start-port=10000 \
  --xray-log-level=none \
  --metrics-protected=true \
  --metrics-username=custom_user \
  --metrics-password=custom_pass
```

### Docker

```bash
docker run -d \
  -e SUBSCRIPTION_URL=https://your-subscription-url/sub \
  -e CHECK_INTERVAL=300 \
  -p 2112:2112 \
  kutovoys/xray-checker
```

### Docker Compose

```yaml
services:
  xray-checker:
    image: kutovoys/xray-checker
    environment:
      - SUBSCRIPTION_URL=https://your-subscription-url/sub
      - CHECK_INTERVAL=300
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
