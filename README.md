# Question Service

API-сервис для работы с вопросами и ответами.

---
## Стек

- **Язык:** Go 1.24
- **HTTP:** `net/http`
- **ORM:** [GORM](https://gorm.io/)
- **База данных:** PostgreSQL 16
- **Миграции:** [pressly/goose](https://github.com/pressly/goose)
- **Логирование:** `uber-go/zap`
- **Сборка/запуск:** Docker, Docker Compose
- **Тесты:** `testing`, `httptest`, `github.com/stretchr/testify`

---
## Структура проекта

```text
question-service/
  cmd/
    api/           # main.go - запуск HTTP API
    migrate/       # main.go - запуск миграций (goose)
  internal/
    app/           # обёртка над http.Server
    config/        # конфиг через env-переменные
    db/            # инициализация GORM + подключение к PostgreSQL
    domain/        # доменные модели Question, Answer
    http/          # HTTP-роутер и хендлеры
    logger/        # обёртка над zap-логгером
    repository/    # интерфейсы и реализации репозиториев на GORM
    service/       # бизнес-логика
    transport/     # общие вспомогательные функции для HTTP-ответов
  migrations/      # SQL-миграции goose
  docker-compose.yml
  go.mod / go.sum
  README.md
```
---
## Запуск приложения

Приложение можно запустить двумя способами: через Docker Compose или локально.

### Способ 1: запуск через Docker Compose 
Выполнить скрипт в директории с файлом docker compose
```bash
docker compose up --build
```
При этом поднимаются три контейнера:

  - **question-service-postgres** -  контейнер с базой данны

  - **question-service-migrate** - контейнер, который выполняет миграции 

  - **question-service-api** - основное приложение (HTTP API)

Все миграции применяются автоматически перед запуском API.

После запуска сервис будет доступен по адресу:
```
http://localhost:8080
```

Проверка состояния:

```bash
GET http://localhost:8080/health
```

---

### Способ 2: локальный запуск (без Docker)
1. Установите все переменные окружения (DB_HOST, DB_USER, DB_PASS, DB_NAME, DB_PORT, DB_SSLMODE).
   2. Запустите сервис:
```bash
   go run cmd/api/main.go
```
3. При необходимости вручную выполнить миграции:
```bash
go run cmd/migrate/main.go
```
После запуска API будет доступен по адресу:
```
http://localhost:8080
```
Проверка состояния:
```bash
GET http://localhost:8080/health
```

---
## HTTP API
Ниже краткое описание основных эндпоинтов.
### Healthcheck
- **GET** /health
  Ответ 200 OK:
```json
{
  "status": "ok"
}
```
### Вопросы
### Создать вопрос
- **POST** /questions
  Пример запроса:

```bash
POST /questions
Content-Type: application/json

{
  "text": "Your question"
}
```

Ответ `201 Created`:

```json
{
  "id": 1,
  "text": "Your question",
  "created_at": "2025-01-01T12:00:00Z"
}
```

### Получить вопрос с ответами

- **GET** `/questions/{id}`

Пример запроса:

```bash
GET /questions/1
```

Ответ `200 OK`:

```json
{
  "id": 1,
  "text": "First question",
  "created_at": "2025-01-01T12:00:00Z",
  "answers": [
    {
      "id": 10,
      "question_id": 1,
      "user_id": "user-123",
      "text": "Example answer",
      "created_at": "2025-01-01T12:10:00Z"
    }
  ]
}
```

### Удалить вопрос

- **DELETE** `/questions/{id}`

Пример запроса:

```bash 
DELETE /questions/1
```

Ответ `204 No Content` — вопрос удалён успешно, и так же все ответы на него


### Ответы

### Создать ответ на вопрос

- **POST** `/questions/{id}/answers`

Тело запроса:

```json
{
  "user_id": "user-123",
  "text": "Answer text"
}

```

DELETE /questions/1

Пример запроса:

```bash
POST /questions/1/answers
Content-Type: application/json

{
  "user_id": "user-123",
  "text": "Answer text"
}
```

Ответ `201 Created`:

```json
{
  "id": 15,
  "question_id": 1,
  "user_id": "user-123",
  "text": "Answer text",
  "created_at": "2025-01-01T13:00:00Z"
}
```

### Получить ответ

- **GET** `/answers/{id}`

Пример запроса:

```bash bash
GET /answers/15
```

Ответ `200 OK`:

```json
{
  "id": 15,
  "question_id": 1,
  "user_id": "user-123",
  "text": "Answer text",
  "created_at": "2025-01-01T13:00:00Z"
}
```
### Удалить ответ

- **DELETE** `/answers/{id}`

Пример запроса:

```bash 
DELETE /answers/1
```

Ответ `204 No Content` — ответ удалён успешно.

---
## Тесты
    Запуск всех тестов с помощью команды go test ./...

## Контакты 
- Почта - max_marinin@mail.ru 
- Telegram - [@mxmrnn](https://t.me/mxmrnn)