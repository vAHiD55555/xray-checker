---
title: GitHub Actions
description: Run Xray Checker using GitHub Actions
---

You can run Xray Checker using GitHub Actions. This approach is useful when you need to run checks from different locations or don't have your own server.

1. Fork the [xray-checker-in-actions](https://github.com/kutovoys/xray-checker-in-actions) repository
2. Configure the following secrets in your forked repository:
   - `SUBSCRIPTION_URL`: Your subscription URL
   - `PUSH_URL`: Prometheus pushgateway URL for metrics collection
   - `INSTANCE`: (Optional) Instance name for metrics identification

The Action will:

- Run every 5 minutes
- Use the latest version of Xray Checker
- Push metrics to your Prometheus pushgateway
- Run with `--run-once` flag to ensure clean execution

This method requires a properly configured Prometheus pushgateway as it can't expose metrics directly. The metrics will be pushed to your specified `PUSH_URL` with the instance label from your configuration.
