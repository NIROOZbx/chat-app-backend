# Chat App - Backend

A real-time chat application backend built with Go, Gin, WebSockets, Redis Pub/Sub, and PostgreSQL.

## Tech Stack

- **Language:** Go 1.25
- **Framework:** Gin
- **Database:** PostgreSQL (Supabase) via pgx/sqlx
- **Cache/Pub-Sub:** Redis
- **Real-time:** WebSockets (gorilla/websocket)
- **Image Upload:** Cloudinary
- **Config:** Viper + YAML + .env
- **Migrations:** golang-migrate
- **Task Runner:** Taskfile

## Project Structure

```
backend/
├── cmd/api/main.go              # Entry point
├── internal/
│   ├── bootstrap/                # App startup and DI wiring
│   ├── handlers/                 # HTTP + WebSocket handlers
│   ├── models/                   # DB models (User, Room, Message)
│   ├── repositories/             # Data access layer
│   ├── routes/                   # Route definitions
│   ├── services/                 # Business logic
│   └── shared/
│       ├── cache/                # Redis client setup
│       ├── config/               # Config loading (Viper)
│       ├── database/             # DB connection
│       ├── dtos/                 # Data transfer objects
│       ├── hub/                  # WebSocket room manager
│       ├── logger/               # Custom logger
│       ├── middleware/           # CORS, logging, session
│       ├── pubsub/               # Redis Pub/Sub + event types
│       ├── request/              # Request DTOs
│       ├── response/             # Response helpers
│       ├── session/              # Redis-backed session store
│       └── utils/                # Image upload, ID parsing, invite codes
├── migrations/                   # SQL migration files
├── config.yaml                   # App configuration
└── Taskfile.yml                  # Task runner commands
```

## Features

- User creation with profile image upload (Cloudinary)
- Session-based auth via Redis (httpOnly cookies)
- Create, join (public/private via invite code), and delete rooms
- Real-time messaging over WebSockets
- Typing indicators via Redis Pub/Sub
- Online user tracking per room
- Paginated message history
- Image upload for rooms and users
- Graceful shutdown
- pprof profiling endpoint (port 6060)

## API Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/` | Health check | No |
| POST | `/api/v1/create-user` | Register user | No |
| POST | `/api/v1/auth-user` | Login (create session) | No |
| GET | `/api/v1/me` | Get current user | Yes |
| POST | `/api/v1/rooms/` | Create room | Yes |
| GET | `/api/v1/rooms/` | List all rooms | No |
| GET | `/api/v1/rooms/joined` | List joined rooms | Yes |
| GET | `/api/v1/rooms/:id` | Get room details | Yes |
| DELETE | `/api/v1/rooms/:id` | Delete room (admin) | Yes |
| POST | `/api/v1/rooms/join/:id` | Join public room | Yes |
| POST | `/api/v1/rooms/join/private/:inviteCode` | Join private room | Yes |
| DELETE | `/api/v1/rooms/leave/:id` | Leave room | Yes |
| GET | `/api/v1/rooms/ws/:id` | WebSocket connection | Yes |
| POST | `/api/v1/rooms/:id/messages` | Send message | Yes |
| GET | `/api/v1/rooms/:id/messages` | Get messages (paginated) | Yes |

## WebSocket Protocol

Connect to `ws://host/api/v1/rooms/ws/:id` with session cookie.

**Client -> Server:**
```json
{ "type": "message.send", "content": "Hello!" }
{ "type": "room.typing" }
```

**Server -> Client (broadcast):**
```json
{ "type": "message.new", "content": "Hello!", "user_name": "john", "created_at": "..." }
{ "type": "user.online", "user_name": "john" }
{ "type": "user.offline", "user_name": "john" }
{ "type": "user.typing", "user_name": "john" }
```

## Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL
- Redis
- Cloudinary account
- [Task](https://taskfile.dev/) (optional)

### Setup

1. Clone the repo and create a `.env` file:

```env
DATABASE_URL=postgresql://user:pass@host:5432/dbname
REDIS_URL=redis://default:pass@host:port
CLOUDINARY_URL=cloudinary://api_key:api_secret@cloud_name
```

2. Update `config.yaml` if needed (defaults to port 8081).

3. Run migrations and start the server:

```sh
task run
```

Or manually:

```sh
task migrate-up
go run cmd/api/main.go
```

### Database Migrations

```sh
task migrate-up            # Apply all pending migrations
task migrate-down 1        # Rollback last migration
task migrate-create my_migration  # Create new migration files
task migrate-force 8       # Force set migration version
```

## Database Schema

- **users** - id, user_name, profile_image, is_online, last_seen, created_at
- **room** - id, name, description, topic, max_members, invite_code, is_private, image, created_at, updated_at
- **room_members** - room_id, user_id, role, is_banned, joined_at (composite PK)
- **messages** - id, room_id, user_id, content, created_at
