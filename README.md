# Refda (رفدة) Backend

Premium Saudi social-fintech API for social gifting and event registries.

## Stack

- **Go 1.23+** with [Gin](https://github.com/gin-gonic/gin)
- **PostgreSQL** + [GORM](https://gorm.io) (data + OTP challenges)
- **JWT** stateless auth
- **go-playground/validator**

## Architecture

```
cmd/api/                 # Entry point
internal/
  app/
    handlers/            # HTTP layer
    services/            # Business logic
    repositories/        # Data access
    dto/                 # Request/response types
    router/              # Route wiring
  domain/                # Entities
  pkg/                   # Shared utilities (config, jwt, otp, storage, …)
migrations/              # SQL reference migrations
postman_collection.json  # API collection
```

## Quick Start

### Prerequisites

- Go 1.21+
- Docker (optional, for Postgres)

### 1. Start infrastructure

```bash
docker compose up -d postgres
```

If the image pull fails with `unexpected EOF`, retry:

```bash
docker compose pull postgres
docker compose up -d postgres
```

Wait until Postgres is healthy (`docker compose ps` shows `healthy`), then run the API.

### 2. Configure

Edit `config.yaml` if needed (defaults match `docker-compose.yml`).

### 3. Run API

```bash
go mod download
make run
```

API: `http://localhost:8080`  
Health: `GET /health`

### OTP (development)

OTP codes are logged to stdout:

```
[MOCK SMS] OTP for +966501234567: 1234
```

Use that code in `POST /api/v1/auth/verify`.

## API Overview

| Group | Endpoint | Auth |
|-------|----------|------|
| Auth | `POST /api/v1/auth/login` | No |
| Auth | `POST /api/v1/auth/register` | No |
| Auth | `POST /api/v1/auth/verify` | No |
| User | `GET /api/v1/users/me` | JWT |
| User | `PUT /api/v1/users/me` | JWT |
| Events | `POST /api/v1/events` | JWT |
| Events | `GET /api/v1/events` | Optional (`?mine=true` + JWT) |
| Events | `GET /api/v1/events/:id` | No |
| Gifts | `POST /api/v1/gifts/contribute` | Optional |
| Gifts | `POST /api/v1/gifts/pay` | No |

Import `postman_collection.json` into Postman for full examples.

Run automated endpoint tests:

```bash
make run   # terminal 1
./scripts/test_api.sh   # terminal 2
```

## Features

- **OTP auth** — Saudi numbers (`+9665XXXXXXXX`), mock SMS
- **JWT** — 7-day default expiry
- **Events** — Wedding, newborn, birthday; public/private; hero image upload
- **Gifts** — Product (group contributions) and cash; atomic `amount_received` updates
- **Payments** — Mock Mada / Apple Pay; reference numbers (`RFD-XXXX`)
- **Localization** — Bilingual error payloads (`message` + `message_ar`)

## Make Commands

```bash
make run          # Run server
make build        # Build binary
make test         # Run tests
make docker-up    # Start all services
make docker-down  # Stop services
```

## License

MIT
