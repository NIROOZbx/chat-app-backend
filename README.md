# Chat App Backend

This is the backend for the Chat App, built with Go. It provides real-time chat functionality using WebSockets, Redis Pub/Sub, and a Postgres database.

## Tech Stack

- **Lanuage:** [Go](https://go.dev/) (v1.25.1+)
- **Framework:** [Gin Web Framework](https://gin-gonic.com/)
- **Database:** [PostgreSQL](https://www.postgresql.org/) with [SQLX](https://github.com/jmoiron/sqlx)
- **Real-time:** [WebSockets](https://github.com/gorilla/websocket) and [Redis Pub/Sub](https://redis.io/)
- **Configuration:** [Viper](https://github.com/spf13/viper) and [Godotenv](https://github.com/joho/godotenv)
- **Storage:** [Cloudinary](https://cloudinary.com/) (for media assets)
- **Task Runner:** [Task](https://taskfile.dev/)

## Prerequisites

Ensure you have the following installed:

- [Go](https://go.dev/doc/install)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Redis](https://redis.io/download/)
- [Task](https://taskfile.dev/installation/) (optional, but recommended)
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) (for database migrations)

## Getting Started

### 1. Clone the repository

```bash
git clone <repository-url>
cd CHAT-APP/backend
```

### 2. Environment Setup

Create a `.env` file in the `backend` directory based on the existing `.env` or configuration requirements.

```env
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/chat_db?sslmode=disable
REDIS_URL=localhost:6379
CLOUDINARY_URL=cloudinary://api_key:api_secret@cloud_name
# Add other necessary variables
```

### 3. Database Migrations

Apply migrations to set up your database schema:

Using Task:
```bash
task migrate-up
```

Or manually:
```bash
migrate -path migrations -database "$DATABASE_URL" up
```

### 4. Running the Application

Using Task:
```bash
task run
```

Or using Go directly:
```bash
go run cmd/api/main.go
```

The server will start on the port specified in your `.env` (default is usually 8080).

## Project Structure

- `cmd/api/`: Application entry point.
- `internal/`: Private application code.
  - `bootstrap/`: Application initialization logic.
  - `handlers/`: HTTP request handlers.
  - `services/`: Business logic and external service integrations (WebSocket, Cloudinary, etc.).
  - `shared/`: Shared utilities, middleware, and configuration.
- `migrations/`: SQL migration files.
- `test/`: Test suites.
- `logs/`: Application logs.

## Available Tasks

| Task | Description |
| :--- | :--- |
| `task run` | Runs database migrations and starts the Go server. |
| `task migrate-create -- <name>` | Creates a new migration file pair. |
| `task migrate-up` | Applies all pending database migrations. |
| `task migrate-down` | Reverts database migrations. |
| `task migrate-force -- <version>` | Forces the migration version. |
