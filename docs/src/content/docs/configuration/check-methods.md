---
title: Check Methods
description: Check methods options and examples
---

Xray Checker supports two proxy checking methods:

## IP method (default)

The `ip` method checks proxy functionality by comparing IP addresses:

1. Gets current IP without using proxy
2. Connects through proxy and gets IP through it
3. Compares IP addresses - if they differ, proxy is considered working

## Status method

The `status` method checks proxy functionality by sending an HTTP request to a specified URL:

1. Connects through proxy
2. Sends GET request to URL specified in `PROXY_STATUS_CHECK_URL` (default is `http://cp.cloudflare.com/generate_204`)
3. If response with code 204 is received, proxy is considered working

This method can be useful when:

- Faster check is needed
- IP method doesn't work due to IP check service blocks
- Need to verify accessibility of specific resource through proxy

To use this method specify:
