---
title: Subscription
description: Subscription options and examples
---

The content of `SUBSCRIPTION_URL` must be in Base64 Encoded format containing a list of proxies. (Standard format for Xray clients - Streisand, V2rayNG).

Proxies with ports 0, 1 will be ignored.

Request headers sent:

```
Accept: */*
User-Agent: Xray-Checker
```
