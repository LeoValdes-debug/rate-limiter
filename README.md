# rate-limiter

HTTP-сервер с rate limiting по IP на основе алгоритма token bucket.

## Стек

- Go 1.22
- Стандартная библиотека (без сторонних зависимостей для core-логики)
- Docker / Docker Compose

## Как работает

Каждый IP получает N токенов (capacity). Один запрос = один токен. Токены пополняются со скоростью refill/сек. Когда токены кончаются - возвращается 429.

Параметры настраиваются через переменные окружения:

| Переменная | По умолчанию | Описание |
|---|---|---|
| RATE_CAPACITY | 10 | максимум токенов на IP |
| RATE_REFILL | 5 | токенов в секунду |
| PORT | 8080 | порт сервера |

## Запуск

```bash
docker-compose up --build
```

## Эндпоинты

```
GET /ping   - статус сервера
GET /hello  - тестовый эндпоинт
```

## Проверить rate limiting

```bash
# отправить 15 запросов подряд, увидеть 429 после 10-го
for i in $(seq 1 15); do curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/hello; done
```

## Тесты

```bash
go test ./internal/limiter/...
```
