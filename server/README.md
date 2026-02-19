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

## Testing

### Test Setup

The test suite includes:
- **Unit Tests**: LWW conflict resolution strategy
- **API Tests**: All HTTP endpoints with mocked dependencies
- **Integration Tests**: Complete sync flow with mocks

### Install Test Dependencies

```bash
# Go modules will automatically install test dependencies
go mod download
go mod tidy
```

The test suite uses:
- `github.com/stretchr/testify` - For assertions and mocks
- `github.com/gin-gonic/gin` - For HTTP testing
- `github.com/google/uuid` - For test UUID generation

### Environment Variables for Database Tests

For tests that require a real database (optional), set these environment variables:

```bash
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5432
export TEST_DB_USER=postgres
export TEST_DB_PASSWORD=postgres
export TEST_DB_NAME=account_test
export TEST_DB_SSLMODE=disable
```

### Running Tests

#### Unit Tests

Run only unit tests (LWW strategy tests):

```bash
cd server
go test -v ./test/unit
```

Or run a specific unit test:

```bash
go test -v ./test/unit -run TestLWWStrategy
```

#### API Tests

Run API endpoint tests (uses mocks, no database required):

```bash
cd server
go test -v ./test/api
```

Or run specific API test suites:

```bash
# Auth API tests
go test -v ./test/api -run TestAuth

# Transaction API tests
go test -v ./test/api -run TestTransaction

# Account API tests
go test -v ./test/api -run TestAccount

# Sync API tests
go test -v ./test/api -run TestPull
```

#### Integration Tests

Run the complete sync flow integration tests:

```bash
cd server
go test -v ./test/integration
```

Or run specific integration tests:

```bash
# Complete sync flow test
go test -v ./test/integration -run TestCompleteSyncFlow

# Multi-device scenario test
go test -v ./test/integration -run TestSyncFlow_MultiDeviceScenario
```

#### Run All Tests

Run the entire test suite:

```bash
cd server
go test -v ./test/...
```

Or with coverage:

```bash
cd server
go test -v -coverprofile=coverage.out ./test/...
go tool cover -html=coverage.out
```

### Test Coverage

To see test coverage:

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./test/...

# View HTML report
go tool cover -html=coverage.out

# View coverage in terminal
go tool cover -func=coverage.out
```

### Test Architecture

#### Directory Structure

```
server/test/
├── testsetup/
│   └── testdb.go           # Test database setup utilities
├── mocks/
│   └── repository_mocks.go # Mock implementations of repositories
├── unit/
│   └── lww_strategy_test.go # LWW conflict resolution unit tests
├── api/
│   ├── auth_api_test.go     # Auth endpoint tests
│   ├── transaction_api_test.go # Transaction endpoint tests
│   ├── account_api_test.go  # Account endpoint tests
│   ├── category_api_test.go # Category endpoint tests
│   └── sync_api_test.go     # Sync endpoint tests
└── integration/
    └── sync_flow_integration_test.go # Complete sync flow tests
```

#### Key Test Files

**Unit Tests (`test/unit/`):**
- `lww_strategy_test.go`: Tests Last-Write-Wins conflict resolution with various scenarios including edge cases.

**API Tests (`test/api/`):**
- `auth_api_test.go`: Tests registration, login, and JWT middleware
- `account_api_test.go`: Tests account CRUD operations
- `transaction_api_test.go`: Tests transaction CRUD and stats endpoints
- `sync_api_test.go`: Tests pull/push sync endpoints

**Integration Tests (`test/integration/`):**
- `sync_flow_integration_test.go`: Tests the complete end-to-end sync flow as documented in CLAUDE.md, including:
  - Device A creates a transaction
  - Device A pushes to server
  - Server notifies Device B via WebSocket
  - Device B pulls changes
  - Conflict resolution scenarios

### Writing New Tests

#### Adding a Unit Test

1. Create file in `test/unit/`
2. Use `github.com/stretchr/testify/assert` for assertions
3. Run with `go test -v ./test/unit`

Example:
```go
package unit

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestYourFeature(t *testing.T) {
    result := YourFunction()
    assert.Equal(t, expected, result)
}
```

#### Adding an API Test

1. Create file in `test/api/`
2. Use `httptest` for HTTP testing
3. Use mocks from `test/mocks/`
4. Setup test router with mock handlers

#### Adding an Integration Test

1. Create file in `test/integration/`
2. Test complete flows across multiple components
3. Use mocks to isolate from external dependencies
4. Follow the sync flow pattern from CLAUDE.md

### Sync Flow Testing

The integration tests follow the exact flow documented in CLAUDE.md:

1. **Handshake**: Client sends device ID and last sync timestamp
2. **Pull**: Server returns changes since last sync
3. **Push**: Client sends local changes
4. **LWW Resolution**: Conflicts resolved by last_modified_at
5. **Commit**: Changes applied to database
6. **Notification**: WebSocket notifies other devices

The integration test `TestCompleteSyncFlow_DeviceACreatesTransaction_DeviceBReceivesSync` verifies this entire flow.

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
