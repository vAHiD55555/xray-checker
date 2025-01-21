---
title: Features
description: Features of Xray Checker
---

## Core Capabilities

### Protocol Support

- VLESS with various security options (TLS, XTLS, Reality)
- VMess with customizable security settings
- Trojan protocol integration
- Shadowsocks with multiple encryption methods

### Monitoring Infrastructure

- Prometheus metrics export
- Individual status endpoints for each proxy
- Latency measurements
- Detailed connection statistics
- Custom instance labeling for distributed setups

### Automation & Management

- Subscription-based configuration management
- Automatic proxy configuration updates
- Dynamic proxy health checking
- Configurable check intervals
- Multiple check methods (IP-based, status-based)

### Security & Integration

- Basic authentication support for sensitive endpoints
- Prometheus pushgateway support
- Uptime Kuma integration
- Flexible deployment options (standalone, Docker, GitHub Actions)

### Web Interface

- Clean, intuitive dashboard
- Real-time status overview
- Configuration details display
- Dark/light theme support
- Mobile-responsive design

## Advanced Features

### Check Methods

- IP-based verification through external services
- Status code verification with customizable endpoints
- Configurable timeouts and retry logic
- Latency simulation options for accurate monitoring

### Deployment Options

- Standalone binary for direct server deployment
- Docker container for containerized environments
- Docker Compose support for orchestrated setups
- GitHub Actions integration for cloud-based monitoring

### Monitoring Capabilities

- Real-time proxy status monitoring
- Latency tracking and reporting
- Success/failure rate metrics
- Instance-based metric segregation
- Custom label support for better organization

### Configuration Management

- Environment variable support
- CLI parameter flexibility
- JSON configuration file support
- Dynamic subscription updates
- Multiple source format support (Base64, JSON, plain text)
