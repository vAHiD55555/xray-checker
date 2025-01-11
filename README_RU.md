# Xray Checker

[![GitHub Release](https://img.shields.io/github/v/release/kutovoys/xray-checker?style=flat&color=blue)](https://github.com/kutovoys/xray-checker/releases/latest)
[![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/kutovoys/xray-checker/build-publish.yml)](https://github.com/kutovoys/xray-checker/actions/workflows/build-publish.yml)
[![DockerHub](https://img.shields.io/badge/DockerHub-kutovoys%2Fxray--checker-blue)](https://hub.docker.com/r/kutovoys/xray-checker/)
[![GitHub License](https://img.shields.io/github/license/kutovoys/xray-checker?color=greeen)](https://github.com/kutovoys/xray-checker/blob/main/LICENSE)
[![ru](https://img.shields.io/badge/lang-ru-blue)](https://github.com/kutovoys/xray-checker/blob/main/README_RU.md)
[![en](https://img.shields.io/badge/lang-en-red)](https://github.com/kutovoys/xray-checker/blob/main/README.md)

Xray Checker - это инструмент для мониторинга доступности прокси-серверов с поддержкой протоколов VLESS, Trojan и Shadowsocks. Он автоматически тестирует соединения через Xray Core и предоставляет метрики для Prometheus, а также API-эндпоинты для интеграции с системами мониторинга.

<div align="center">
  <img src=".github/screen/xray-checker.png" alt="Dashboard Screenshot">
</div>

## Возможности

- **Поддержка протоколов**: VLESS, Trojan и Shadowsocks
- **Метрики Prometheus**: Экспорт метрик состояния прокси для Prometheus
- **API-эндпоинты**: Отдельные эндпоинты для каждого прокси для интеграции с системой мониторинга
- **Автоматические обновления**: Периодическая проверка и обновление конфигурации из URL подписки
- **Веб-интерфейс**: Простой интерфейс для просмотра статуса прокси и конфигурации
- **Базовая аутентификация**: Опциональная защита метрик и API с помощью базовой аутентификации
- **Поддержка Docker**: Простое развертывание с использованием Docker и Docker Compose

## Метрики

Экспортер предоставляет следующие метрики:

| Название                | Описание                                    |
| ----------------------- | ------------------------------------------- |
| `xray_proxy_status`     | Статус прокси (1: работает, 0: не работает) |
| `xray_proxy_latency_ms` | Задержка через прокси в миллисекундах       |

Каждая метрика включает следующие лейблы:

- `protocol`: Тип протокола (vless/trojan/shadowsocks)
- `address`: Адрес сервера и порт
- `name`: Имя конфигурации прокси
- `instance`: Опциональная метка instance (если указана через --metrics-instance)

## Конфигурация

| Переменная окружения           | Аргумент командной строки        | По умолчанию                            | Описание                                            |
| ------------------------------ | -------------------------------- | --------------------------------------- | --------------------------------------------------- |
| **Subscription**               |
| `SUBSCRIPTION_URL`             | `--subscription-url`             | -                                       | URL подписки для получения конфигураций             |
| `SUBSCRIPTION_UPDATE`          | `--subscription-update`          | `true`                                  | Обновлять подписку автоматически                    |
| `SUBSCRIPTION_UPDATE_INTERVAL` | `--subscription-update-interval` | `300`                                   | Интервал обновления подписки в секундах             |
| **Proxy**                      |
| `PROXY_CHECK_INTERVAL`         | `--proxy-check-interval`         | `300`                                   | Интервал проверки в секундах                        |
| `PROXY_CHECK_METHOD`           | `--proxy-check-method`           | `ip`                                    | Метод проверки (ip/status)                          |
| `PROXY_IP_CHECK_URL`           | `--proxy-ip-check-url`           | `https://api.ipify.org?format=text`     | URL сервиса проверки IP                             |
| `PROXY_STATUS_CHECK_URL`       | `--proxy-status-check-url`       | `http://cp.cloudflare.com/generate_204` | URL для проверки статуса                            |
| `PROXY_TIMEOUT`                | `--proxy-timeout`                | `30`                                    | Таймаут проверки в секундах                         |
| `SIMULATE_LATENCY`             | `--simulate-latency`             | `true`                                  | Добавлять задержку к ответу                         |
| **Xray**                       |
| `XRAY_START_PORT`              | `--xray-start-port`              | `10000`                                 | Начальный порт для конфигураций                     |
| `XRAY_LOG_LEVEL`               | `--xray-log-level`               | `none`                                  | Уровень логирования (debug/info/warning/error/none) |
| **Metrics**                    |
| `METRICS_PORT`                 | `--metrics-port`                 | `2112`                                  | Порт для метрик                                     |
| `METRICS_PROTECTED`            | `--metrics-protected`            | `false`                                 | Защита метрик Basic Auth                            |
| `METRICS_USERNAME`             | `--metrics-username`             | `metricsUser`                           | Имя пользователя для Basic Auth                     |
| `METRICS_PASSWORD`             | `--metrics-password`             | `MetricsVeryHardPassword`               | Пароль для Basic Auth                               |
| `METRICS_PUSH_URL`             | `--metrics-push-url`             | -                                       | URL Prometheus pushgateway                          |
| `METRICS_INSTANCE`             | `--metrics-instance`             | -                                       | Метка instance для метрик                           |
| **Другое**                     |
| `RUN_ONCE`                     | `--run-once`                     | `false`                                 | Выполнение одиного цикла проверки и выход           |

### Формат подписки

Содержимое `SUBSCRIPTION_URL` должно быть в формате Base64 Encoded списка прокси. (Стандартный формат Xray-клиентов – Streisand, V2rayNG).

Прокси с портами 0, 1 – будут игнорироваться.

Отправляемые заголовки:

```
Accept: */*
User-Agent: Xray-Checker
```

## Использование

### CLI

```bash
# Минимально необходимая конфигурация
./xray-checker --subscription-url="https://your-subscription-url/sub"
```

```bash
# Использование всех доступных параметров
./xray-checker \
  --subscription-url="https://your-subscription-url/sub" \
  --subscription-update=true \
  --subscription-update-interval=300 \
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
  --metrics-password=custom_pass \
  --metrics-instance=node-1 \
  --metrics-push-url="https://push.example.com" \
  --run-once=false
```

### Docker

```bash
# Минимально необходимая конфигурация
docker run -d \
  -e SUBSCRIPTION_URL=https://your-subscription-url/sub \
  -p 2112:2112 \
  kutovoys/xray-checker
```

```bash
# Использование всех доступных параметров
docker run -d \
  -e SUBSCRIPTION_URL=https://your-subscription-url/sub \
  -e SUBSCRIPTION_UPDATE=true \
  -e SUBSCRIPTION_UPDATE_INTERVAL=300 \
  -e PROXY_CHECK_INTERVAL=300 \
  -e PROXY_CHECK_METHOD=ip \
  -e PROXY_TIMEOUT=30 \
  -e PROXY_IP_CHECK_URL=https://api.ipify.org?format=text \
  -e PROXY_STATUS_CHECK_URL=http://cp.cloudflare.com/generate_204 \
  -e SIMULATE_LATENCY=true \
  -e XRAY_START_PORT=10000 \
  -e XRAY_LOG_LEVEL=none \
  -e METRICS_PORT=2112 \
  -e METRICS_PROTECTED=true \
  -e METRICS_USERNAME=custom_user \
  -e METRICS_PASSWORD=custom_pass \
  -e METRICS_INSTANCE=node-1 \
  -e METRICS_PUSH_URL=https://push.example.com \
  -e RUN_ONCE=false \
  -p 2112:2112 \
  kutovoys/xray-checker
```

### Docker Compose

```yaml
# Минимальная конфигурация
services:
  xray-checker:
    image: kutovoys/xray-checker
    environment:
      - SUBSCRIPTION_URL=https://your-subscription-url/sub
    ports:
      - "2112:2112"
```

```yaml
# Полная конфигурация со всеми доступными параметрами
services:
  xray-checker:
    image: kutovoys/xray-checker
    environment:
      - SUBSCRIPTION_URL=https://your-subscription-url/sub
      - SUBSCRIPTION_UPDATE=true
      - SUBSCRIPTION_UPDATE_INTERVAL=300
      - PROXY_CHECK_INTERVAL=300
      - PROXY_CHECK_METHOD=ip
      - PROXY_TIMEOUT=30
      - PROXY_IP_CHECK_URL=https://api.ipify.org?format=text
      - PROXY_STATUS_CHECK_URL=http://cp.cloudflare.com/generate_204
      - SIMULATE_LATENCY=true
      - XRAY_START_PORT=10000
      - XRAY_LOG_LEVEL=none
      - METRICS_PORT=2112
      - METRICS_PROTECTED=true
      - METRICS_USERNAME=custom_user
      - METRICS_PASSWORD=custom_pass
      - METRICS_INSTANCE=node-1
      - METRICS_PUSH_URL=https://push.example.com
      - RUN_ONCE=false
    ports:
      - "2112:2112"
```

### GitHub Actions

Вы можете запускать Xray Checker с помощью GitHub Actions. Этот подход полезен, когда нужно выполнять проверки с разных локаций или у вас нет собственного сервера.

1. Сделайте форк репозитория [xray-checker-in-actions](https://github.com/kutovoys/xray-checker-in-actions)
2. Настройте следующие секреты в вашем форке:
   - `SUBSCRIPTION_URL`: URL вашей подписки
   - `PUSH_URL`: URL Prometheus pushgateway для сбора метрик
   - `INSTANCE`: (Опционально) Имя инстанса для идентификации метрик

Action будет:

- Запускаться каждые 5 минут
- Использовать последнюю версию Xray Checker
- Отправлять метрики в ваш Prometheus pushgateway
- Запускаться с флагом `--run-once` для обеспечения чистого выполнения

Этот метод требует правильно настроенного Prometheus pushgateway, так как не может напрямую экспортировать метрики. Метрики будут отправляться на указанный вами `PUSH_URL` с меткой instance из вашей конфигурации.

### Конфигурация Prometheus

Добавьте следующее в ваш prometheus.yml:

```yaml
scrape_configs:
  - job_name: "xray-checker"
    static_configs:
      - targets: ["localhost:2112"]
    scrape_interval: 1m
```

## API эндпоинты

- `/` - Информационная страница
- `/metrics` - Эндпоинт метрик Prometheus
- `/health` - Эндпоинт проверки работоспособности
- `/config/{protocol}-{address}-{port}` - Статус конкретного прокси (возвращает 200 OK, если работает, 503, если не работает)

### Интеграция с Uptime Kuma

Вы можете отслеживать каждый прокси с помощью его отдельного эндпоинта в Uptime Kuma:

1. Добавьте новый монитор
2. Выберите "HTTP(s)"
3. Введите URL: `http://your-server:2112/config/vless-example.com-443`
4. Монитор покажет "Up", если прокси работает, и "Down", если он не работает

## Логика проверки соединения

1. Начальная настройка:

   - Получение конфигураций из URL подписки
   - Генерация унифицированного файла конфигурации Xray
   - Запуск экземпляра Xray Core

2. Периодические проверки (каждые N минут):
   - Получение текущего IP без прокси
   - Для каждой конфигурации прокси:
     - Подключение через локальный SOCKS5 порт
     - Попытка получить IP через прокси
     - Сравнение IP для определения работоспособности прокси
     - Обновление метрик Prometheus и внутреннего статуса
   - Проверка URL подписки на изменения
     - При обнаружении изменений:
       - Генерация новой конфигурации Xray
       - Перезапуск экземпляра Xray Core
       - Обновление эндпоинтов

## Участие в разработке

Мы рады любому вкладу в развитие Xray Checker. Будь то исправление ошибок, новые функции или улучшение документации, вот краткое руководство по участию:

1. **Fork & Branch**: Сделайте форк этого репозитория и создайте ветку для ваших изменений
2. **Implement**: Внесите изменения, сохраняя код чистым и документированным
3. **Test**: Убедитесь, что ваши изменения не нарушают существующую функциональность
4. **Commit & PR**: Создавайте коммиты с четкими сообщениями и откройте Pull Request
5. **Feedback**: Будьте готовы к обсуждению и улучшению вашего вклада

Если вы новичок, GitHub предоставляет отличный гайд по [созданию pull request](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request).
