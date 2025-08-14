# E-Commerce Analytics & Personalization System

A comprehensive Phase 1 implementation of an e-commerce analytics and personalization system built with Go, featuring real-time event processing, rule-based user segmentation, and automated offer generation.

## Features

### 🚀 Phase 1 Implementation
- **Event Simulation & Integration**: HTTP endpoints to receive e-commerce events (page views, purchases, cart actions)
- **Rule-Based Analysis**: Intelligent user segmentation using configurable rules
- **Real-time Segmentation**: Dynamic user classification (High Spender, Window Shopper, Discount Seeker)
- **Automated Offer Generation**: Personalized coupons and discounts based on user segments
- **WebSocket Support**: Real-time notifications for events, segment changes, and offer generation
- **E-commerce Simulation**: Shopify webhook simulation for testing

### 🛠️ Technical Stack
- **Language**: Go 1.23+
- **Framework**: Gorilla Mux (HTTP routing)
- **Database**: PostgreSQL with raw SQL
- **Cache**: Redis for session & segment lookups
- **WebSockets**: Gorilla WebSocket for real-time communication
- **Authentication**: JWT-based authentication
- **Containerization**: Docker & Docker Compose
- **Dependency Injection**: Google Wire

## 🏗️ Architecture

```
├── cmd/api/           # Application entry point
├── internal/          # Private application code
│   ├── config/        # Configuration management
│   ├── controllers/   # HTTP handlers
│   ├── models/        # Data models
│   ├── repositories/  # Data access layer
│   ├── services/      # Business logic
│   └── utils/         # Utility functions
├── pkg/               # Public packages
├── scripts/           # Database migrations
└── web/               # Static files (if any)
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL
- Redis
- Make (optional)

### Setup

1. **Clone and install dependencies:**
   ```bash
   git clone <repository>
   cd <repository>
   go mod download
   ```

2. **Configure environment:**
   ```bash
   cp config.yaml.example config.yaml
   # Edit config.yaml with your database and Redis settings
   ```

3. **Run database migrations:**
   ```bash
   make migrate-up
   ```

4. **Start the server:**
   ```bash
   make run
   ```

The API will be available at `http://localhost:8080`

## 📋 API Endpoints

### Users
- `GET /api/v1/users` - Get all users
- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

### Events
- `POST /api/v1/events` - Process event batch
- `GET /api/v1/events/{user_id}` - Get events for user

### Health
- `GET /health` - Health check

## 🛠️ Development

### Available Make Commands
```bash
make run          # Run the application
make build        # Build the application
make test         # Run tests
make migrate-up   # Run database migrations
make migrate-down # Rollback database migrations
make dev-setup    # Setup development environment
```

### Project Structure

#### Configuration (`internal/config/`)
- `config.go` - Main configuration structure
- `database.go` - Database configuration
- `redis.go` - Redis configuration

#### Models (`internal/models/`)
- `user.go` - User model
- `event.go` - Event models

#### Controllers (`internal/controllers/`)
- `user_controller.go` - User HTTP handlers
- `event_controller.go` - Event HTTP handlers

#### Services (`internal/services/`)
- `user_service.go` - User business logic
- `event_service.go` - Event business logic

#### Repositories (`internal/repositories/`)
- `user_repository.go` - User data access

## 🔧 Configuration

The application uses Viper for configuration management. Key configuration options:

```yaml
server:
  port: 8080
  host: "0.0.0.0"

database:
  host: "localhost"
  port: 5432
  name: "app_db"
  user: "postgres"
  password: "password"
  ssl_mode: "disable"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
```

## 📊 Database Schema

### Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### User Events Table
```sql
CREATE TABLE user_events (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    type VARCHAR(32) NOT NULL,
    data TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/services
```

## 🚀 Deployment

### Docker
```bash
# Build image
docker build -t app .

# Run container
docker run -p 8080:8080 app
```

### Docker Compose
```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down
```

## 📝 License

This project is licensed under the MIT License. 