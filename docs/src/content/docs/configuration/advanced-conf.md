---
title: Check Methods
description: Advanced configuration options
---

## Advanced Configuration

### Custom IP Check Services

You can use alternative IP check services:

- `http://ip.sb`
- `https://api64.ipify.org`
- `http://ifconfig.me`

Example:

```bash
PROXY_IP_CHECK_URL=http://ip.sb
```

### Custom Status Check URLs

Alternative status check URLs:

- `http://www.gstatic.com/generate_204`
- `http://www.qualcomm.cn/generate_204`
- `http://cp.cloudflare.com/generate_204`

Example:

```bash
PROXY_STATUS_CHECK_URL=http://www.gstatic.com/generate_204
```

### Security Configuration

Enable authentication for sensitive endpoints:

```bash
METRICS_PROTECTED=true
METRICS_USERNAME=custom_user
METRICS_PASSWORD=secure_password
```

### Instance Labeling

Add instance labels for distributed setups:

```bash
METRICS_INSTANCE=datacenter-1
```

### Update Intervals

Customize check and update intervals:

```bash
# Check every minute
PROXY_CHECK_INTERVAL=60

# Update subscription every hour
SUBSCRIPTION_UPDATE_INTERVAL=3600
```

### Logging Configuration

Adjust Xray Core logging:

```bash
# Enable debug logging
XRAY_LOG_LEVEL=debug

# Disable logging
XRAY_LOG_LEVEL=none
```

### Port Configuration

Customize port ranges:

```bash
# Start SOCKS5 ports from 20000
XRAY_START_PORT=20000

# Change metrics port
METRICS_PORT=9090
```
