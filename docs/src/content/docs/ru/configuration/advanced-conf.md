---
title: Расширенная конфигурация
description: Расширенные параметры конфигурации
---

### Альтернативные сервисы проверки IP

Вы можете использовать альтернативные сервисы проверки IP (подробнее в разделе [методы проверки](/ru/configuration/check-methods)):

- `http://ip.sb`
- `https://api64.ipify.org`
- `http://ifconfig.me`

Пример:

```bash
PROXY_IP_CHECK_URL=http://ip.sb
```

### Альтернативные URL проверки статуса

Альтернативные URL для проверки статуса (подробнее в разделе [методы проверки](/ru/configuration/check-methods)):

- `http://www.gstatic.com/generate_204`
- `http://www.qualcomm.cn/generate_204`
- `http://cp.cloudflare.com/generate_204`

Пример:

```bash
PROXY_STATUS_CHECK_URL=http://www.gstatic.com/generate_204
```

### Настройка безопасности

Включение аутентификации для чувствительных эндпоинтов:

```bash
METRICS_PROTECTED=true
METRICS_USERNAME=custom_user
METRICS_PASSWORD=secure_password
```

### Маркировка экземпляров

Добавление меток экземпляров для распределенных установок:

```bash
METRICS_INSTANCE=datacenter-1
```

### Интервалы обновления

Настройка интервалов проверки и обновления:

```bash
# Проверка каждую минуту
PROXY_CHECK_INTERVAL=60

# Обновление подписки каждый час
SUBSCRIPTION_UPDATE_INTERVAL=3600
```

### Настройка логирования

Настройка логирования Xray Core:

```bash
# Включение отладочного логирования
XRAY_LOG_LEVEL=debug

# Отключение логирования
XRAY_LOG_LEVEL=none
```

### Настройка портов

Настройка диапазонов портов:

```bash
# Начало портов SOCKS5 с 20000
XRAY_START_PORT=20000

# Изменение порта метрик
METRICS_PORT=9090
```
