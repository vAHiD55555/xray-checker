---
title: Check Methods
description: Check methods options and examples
---

## Check Methods

Xray Checker supports two methods for verifying proxy functionality:

### IP Check Method (Default)

```bash
--proxy-check-method=ip
```

This method:

1. Gets current IP without proxy
2. Connects through proxy
3. Gets IP through proxy
4. Compares IPs to verify proxy is working

Benefits:

- More reliable verification
- Confirms actual proxy functionality
- Detects transparent proxies

Configuration:

```bash
PROXY_CHECK_METHOD=ip
PROXY_IP_CHECK_URL=https://api.ipify.org?format=text
PROXY_TIMEOUT=30
```

### Status Check Method

```bash
--proxy-check-method=status
```

This method:

1. Connects through proxy
2. Requests specified URL
3. Verifies response status code

Benefits:

- Faster verification
- Lower bandwidth usage
- Works with restrictive firewalls

Configuration:

```bash
PROXY_CHECK_METHOD=status
PROXY_STATUS_CHECK_URL=http://cp.cloudflare.com/generate_204
PROXY_TIMEOUT=30
```
