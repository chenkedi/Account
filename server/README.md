# Account Server

Go server for the personal financial management system.

## Development Setup

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Redis 7+

### Quick Start

1. Copy the example config:
   ```bash
   cp config.yaml.example config.yaml
   ```

2. Update config.yaml with your database credentials

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Start PostgreSQL and Redis

5. Run the server:
   ```bash
   go run cmd/server/main.go
   ```

## Project Structure

```
server/
├── cmd/
│   └── server/
│       └── main.go          # Entry point
├── internal/
│   ├── api/
│   │   ├── handlers/         # HTTP handlers
│   │   ├── middleware/       # Gin middleware
│   │   └── routes.go         # Route definitions
│   ├── business/
│   │   ├── models/           # Domain models
│   │   └── services/         # Business logic
│   ├── data/
│   │   ├── repository/       # Data access layer
│   │   └── database/
│   │       ├── migrations/   # SQL migrations
│   │       └── postgres.go   # DB connections
│   └── sync/                 # Sync engine
├── pkg/
│   ├── auth/                 # JWT auth
│   ├── config/               # Config management
│   ├── logger/               # Logging
│   └── utils/                # Utilities
└── api/                      # API specs
```

## API Endpoints

### Public Endpoints
- `GET /health` - Health check
- `POST /api/v1/auth/register` - User registration (creates default categories)
- `POST /api/v1/auth/login` - User login
- `GET /ws/sync?token=<token>&device_id=<device_id>` - WebSocket real-time sync notifications

### Protected Endpoints (require JWT in Authorization header)
- `GET /api/v1/me` - Get current user info

### Accounts
- `POST /api/v1/accounts` - Create account
- `GET /api/v1/accounts` - Get all accounts
- `GET /api/v1/accounts/:id` - Get account by ID
- `PUT /api/v1/accounts/:id` - Update account
- `DELETE /api/v1/accounts/:id` - Delete account

### Categories
- `POST /api/v1/categories` - Create category
- `GET /api/v1/categories` - Get all categories
- `GET /api/v1/categories/type/:type` - Get categories by type (income/expense)
- `GET /api/v1/categories/:id` - Get category by ID
- `PUT /api/v1/categories/:id` - Update category
- `DELETE /api/v1/categories/:id` - Delete category

### Transactions
- `POST /api/v1/transactions` - Create transaction (updates account balance)
- `GET /api/v1/transactions` - Get transactions (paginated)
- `GET /api/v1/transactions/range?start_date=&end_date=` - Get transactions by date range
- `GET /api/v1/transactions/stats?start_date=&end_date=` - Get stats for date range
- `GET /api/v1/transactions/:id` - Get transaction by ID
- `PUT /api/v1/transactions/:id` - Update transaction
- `DELETE /api/v1/transactions/:id` - Delete transaction

### Sync
- `POST /api/v1/sync/pull` - Get server changes since last sync
- `POST /api/v1/sync/push` - Push local changes to server (LWW conflict resolution)

## Sync Protocol

The sync engine uses Last-Write-Wins (LWW) conflict resolution based on `last_modified_at` timestamps.

### Sync Flow:
1. **Pull**: Request all changes from server since last sync
2. **Push**: Send all local changes to server
3. **Resolve**: Server applies LWW strategy to resolve any conflicts
4. **Notify**: WebSocket notifies other devices to sync

### WebSocket Messages:
- `sync_available`: Notification that new changes are available to pull

All entities include `last_modified_at`, `version`, and `is_deleted` for soft delete support.

## Database

Migrations are automatically run on server startup. To manage migrations manually:

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path internal/data/database/migrations -database "postgres://..." up

# Rollback migrations
migrate -path internal/data/database/migrations -database "postgres://..." down 1
```
