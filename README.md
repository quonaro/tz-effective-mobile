# Subscriptions Manager

Fullstack приложение для управления онлайн подписками. Backend на Go (Huma + CHI), Frontend на Nuxt 4 + shadcn/ui.

**🚀 Тестовый сервер:** https://effective-mobile.tz.d3vb0x.ru

**📚 Документация API (Swagger UI):** https://effective-mobile.tz.d3vb0x.ru/api/docs
**📝 Technical Specification:** https://effective-mobile.tz.d3vb0x.ru/tz


## Структура проекта

```
.
├── backend/          # Go backend (Huma + CHI + PostgreSQL)
├── frontend/         # Nuxt 4 frontend (shadcn/ui)
├── docker-compose.yml
├── .env
└── .gitignore
```

## Требования

- Docker и Docker Compose
- Go 1.23+ (для локальной разработки backend)
- Node.js 20+ и pnpm 9+ (для локальной разработки frontend)

## Запуск через Docker Compose

```bash
# Скопировать .env.example в .env (если ещё нет)
cp .env.example .env

# Запустить все сервисы
docker compose up --build

# Backend: http://localhost:8080
# Frontend: http://localhost:3000
# Swagger: http://localhost:8080/docs
```

## Локальная разработка

### Backend

```bash
cd backend

# Установка зависимостей
go mod download

# Запуск (требуется PostgreSQL)
DATABASE_URL=postgres://subs:subs@localhost:5432/subscriptions?sslmode=disable go run ./cmd/subscriptions

# Линтер
golangci-lint run ./...
```

### Frontend

```bash
cd frontend

# Установка зависимостей
pnpm install

# Запуск dev-сервера
pnpm dev

# Сборка
pnpm build

# Линтеры
pnpm lint:oxlint
pnpm lint:prettier
pnpm format
```

## API Endpoints

Базовый URL: `http://localhost:8080`

Swagger UI: `http://localhost:8080/docs`
OpenAPI spec: `http://localhost:8080/openapi.json`

### Создать подписку

`POST /subscriptions`

**Тело запроса:**
```json
{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": "12-2025"
}
```

**Ответ:** `201 Created`
```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "2025-07-01T00:00:00Z",
  "end_date": "2025-12-01T00:00:00Z",
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z"
}
```

### Получить подписку по ID

`GET /subscriptions/{id}`

**Параметры:**
- `id` — UUID подписки (path parameter)

**Ответ:** `200 OK` или `404 Not Found`

### Список подписок

`GET /subscriptions`

**Query-параметры:**
- `user_id` — фильтр по пользователю (UUID)
- `service_name` — фильтр по названию сервиса
- `limit` — лимит записей (по умолчанию 100)
- `offset` — смещение для пагинации

**Пример:** `GET /subscriptions?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&limit=10`

**Ответ:**
```json
{
  "data": [...],
  "total": 42
}
```

### Обновить подписку

`PUT /subscriptions/{id}`

**Тело запроса:** (все поля опциональны)
```json
{
  "service_name": "Yandex Plus Pro",
  "price": 500,
  "start_date": "08-2025",
  "end_date": "11-2025"
}
```

**Ответ:** `200 OK`

### Удалить подписку

`DELETE /subscriptions/{id}`

**Ответ:** `200 OK`

### Подсчитать суммарную стоимость

`GET /subscriptions/total`

**Query-параметры:**
- `user_id` — фильтр по пользователю (UUID, опционально)
- `service_name` — фильтр по сервису (опционально)
- `start_month` — начало периода, формат `MM-YYYY` (обязательно)
- `end_month` — конец периода, формат `MM-YYYY` (обязательно)

**Пример:** `GET /subscriptions/total?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&start_month=01-2025&end_month=12-2025`

**Ответ:**
```json
{
  "total_cost": 4800
}
```

### Формат дат

- `start_date`, `end_date` в теле запроса — `MM-YYYY` (например, `07-2025`)
- В ответах ISO 8601 (например, `2025-07-01T00:00:00Z`)

## Миграции

### Запуск миграций

```bash
cd backend

# Применить все миграции
go run ./cmd/migrate -direction=up

# Откатить миграции
go run ./cmd/migrate -direction=down
```

### Генерация тестовых данных (1000 записей)

Миграция `001_seed_data.sql` генерирует 1000 случайных подписок:
- 10 различных пользователей (UUID)
- 150+ различных сервисов (Netflix, Spotify, и т.д.)
- Случайная цена от 100 до 10000
- Случайные даты подписки за последние 2 года
- 30% записей имеют дату окончания

## Тестовые данные

UUID пользователей для тестирования: см. [`docs/UUID`](./docs/UUID)

## Pre-commit hooks

Установить pre-commit:

```bash
pip install pre-commit
pre-commit install
```

Хуки запускают:
- golangci-lint для Go
- prettier для frontend
- oxlint для frontend

## Technical Specification

Full technical specification is available on the Spec page: `/tz`

**Local:** http://localhost:3000/tz  
**Test Server:** https://effective-mobile.tz.d3vb0x.ru/tz

The specification file is located at [`docs/ТЗ.md`](./docs/ТЗ.md).
