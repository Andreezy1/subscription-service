# Subscription Service

![Go](https://img.shields.io/badge/Go-1.26-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-blue)
![Docker](https://img.shields.io/badge/Docker-ready-blue)

REST API для управления пользовательскими подписками.

Сервис позволяет создавать, получать, изменять и удалять подписки пользователей, а также рассчитывать суммарную стоимость подписок за указанный период с возможностью фильтрации по пользователю и названию сервиса.

---

## Возможности

* Создание подписки
* Получение списка подписок с фильтрацией
* Получение подписки по идентификатору
* Обновление подписки
* Удаление подписки
* Расчет суммарной стоимости подписок за период
* Валидация входных данных
* Логирование HTTP-запросов и ошибок
* Swagger-документация
* Развертывание через Docker и Docker Compose
* Автоматическое применение миграций базы данных при старте с помощью контейнера
* Graceful Shutdown (корректное завершение работы)
* Unit-тестирование слоев приложения

---

## Используемые технологии

* **Language:** Go
* **Router:** Chi Router (v5)
* **Database:** PostgreSQL (драйвер pgx)
* **Logging:** slog (структурированное логирование)
* **Documentation:** Swagger (swaggo)
* **Containerization:** Docker, Docker Compose
* **Migrations:** golang-migrate (отдельный Docker-контейнер)
* **Testing:** Go testing package

---

## Архитектура проекта

Проект построен по принципам многослойной (Layered) архитектуры.
Зависимости передаются через конструкторы (Dependency Injection), благодаря чему слои остаются независимыми и легко тестируются.

```
HTTP Request
│
Handler  <─── (Обработка HTTP запросов и логирование запросов)
│
Service  <─── (Бизнес-логика)
│
Repository <─── (Прямая работа с БД, SQL-запросы)
│
PostgreSQL
```

### Handler

Отвечает за:

* обработку HTTP-запросов;
* разбор JSON;
* валидацию параметров запроса;
* преобразование DTO ↔ Model;
* формирование HTTP-ответов.

### Service

Содержит бизнес-логику приложения.

### Repository

Работает с PostgreSQL.

Содержит SQL-запросы и преобразование результатов в модели.

---

## Структура проекта

```
.
├── cmd/
│   └── app/
│       ├── app.go  # точка входа приложения
│       ├── logger.go
│       └── main.go
│
├── config/         # загрузка конфигурации
│
├── docs/           # Swagger
│
├── internal/
│   ├── handler/    # HTTP handlers
│   ├── model/      # модели и ошибки
│   ├── repository/ # работа с PostgreSQL
│   └── service/    # бизнес-логика
│
├── migrations/     # SQL-миграции
│
├── .env
├── Dockerfile
├── docker-compose.yml
└── README.md
```

---
## Требования

Для запуска необходимы:

- Docker
- Docker Compose

## Запуск проекта

### Клонировать репозиторий

```bash
git clone <repository_url>

cd subscription-service
```

---

### Создать файл .env

```
HTTP_PORT=8080

DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscriptions
DB_SSLMODE=disable
```

---

### Запустить проект

```bash
docker compose up --build
```

При этом Docker Compose автоматически:

   * Поднимет PostgreSQL и дождется его готовности (healthcheck).

   * Запустит контейнер golang-migrate и накатит все миграции из папки /migrations.

   * Скомпилирует и запустит Go-сервис.

---

## Swagger

После запуска документация доступна по адресу:

```
http://localhost:8080/swagger/index.html
```

Через Swagger UI можно выполнять запросы к API.

---

## API

### Формат дат

Во всех запросах используется формат:

MM-YYYY

Например

01-2025
12-2026

### Создать подписку

```
URL: POST /subscriptions
```

Пример запроса

```json
{
  "service_name": "Netflix",
  "price": 799,
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "start_date": "01-2025",
  "end_date": "12-2025"
}
```

Пример успешного ответа (201 Created):

```json
{
  "id": 42,
  "service_name": "Netflix",
  "price": 799,
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "start_date": "01-2025",
  "end_date": "12-2025"
}
```
---

### Получить список подписок
```
URL: GET /subscriptions
```

Необязательные query-параметры

| Параметр | Описание |
|----------|----------|
| user_id | UUID пользователя |
| service_name | Название сервиса |

Пример
GET /subscriptions?service_name=Netflix

Пример ответа

```json
[
  {
    "id": 1,
    "service_name": "Netflix",
    "price": 799,
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "start_date": "01-2025",
    "end_date": "12-2025"
  }
]
```

### Получить подписку

```
URL: GET /subscriptions/{id}
```
Пример успешного ответа (200 OK):
```json
{
  "id": 42,
  "service_name": "Netflix",
  "price": 799,
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "start_date": "01-2025",
  "end_date": "12-2025"
}
```
---

### Обновить подписку

```
URL: PUT /subscriptions/{id}
```

---

### Удалить подписку

```
URL: DELETE /subscriptions/{id}
```

---

### Рассчитать стоимость подписок
Возвращает общую сумму затрат на подписки за указанный период с возможностью фильтрации.

```
URL: GET /subscriptions/total
```

Параметры:

| Параметр     | Обязательный | Описание               |
| ------------ | ------------ | ---------------------- |
| start_date   | Да           | начало периода         |
| end_date     | Да           | конец периода          |
| user_id      | Нет          | фильтр по пользователю |
| service_name | Нет          | фильтр по сервису      |

Пример

```
GET /subscriptions/total?start_date=01-2025&end_date=12-2025&user_id=550e8400-e29b-41d4-a716-446655440000
```

---

## Логирование

Используется стандартный пакет Go `log/slog`.

Логируются:

- входящие HTTP-запросы;
- request_id;
- HTTP-статус ответа;
- время обработки запроса;
- ошибки приложения;
- успешные бизнес-операции.

Пример

```
level=INFO
msg="http request"
request_id=...
method=POST
path=/subscriptions
status=201
duration=2ms
```

---

## Миграции

Для управления схемой базы данных используются SQL-миграции.

При запуске Docker Compose миграции применяются автоматически.

---

## Тестирование

Запуск всех тестов

```bash
go test ./...
```

Покрыты тестами:

- model
- service
- handler (парсинг параметров и обработка запросов)

---

## Graceful Shutdown

Приложение корректно завершает работу при получении SIGINT/SIGTERM:

* прекращает принимать новые соединения;
* ожидает завершения активных запросов;
* закрывает HTTP-сервер;
* закрывает соединение с PostgreSQL.

---

## Особенности реализации

* многослойная архитектура;
* внедрение зависимостей через конструкторы (Dependency Injection);
* централизованная обработка ошибок;
* единое логирование через slog;
* middleware для логирования HTTP-запросов;
* конфигурация через переменные окружения;
* Docker-ready приложение;
* автоматические миграции;
* Swagger-документация;
* unit-тесты.

---

## Использованные практики

В проекте применены:

- Layered Architecture
- Dependency Injection
- REST API
- Repository Pattern
- Structured Logging
- Graceful Shutdown
- Docker Compose
- Database Migrations
- Swagger/OpenAPI
- Unit Testing
