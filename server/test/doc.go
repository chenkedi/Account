/*
Package test contains the complete test suite for the account server.

Test Structure:

  test/
  ├── testsetup/          # Test utilities and setup
  │   └── testdb.go       # Test database configuration
  ├── mocks/              # Mock implementations
  │   └── repository_mocks.go  # Mock repositories for testing
  ├── unit/               # Unit tests
  │   └── lww_strategy_test.go # LWW conflict resolution tests
  ├── api/                # API endpoint tests
  │   ├── auth_api_test.go      # Auth endpoint tests
  │   ├── account_api_test.go   # Account endpoint tests
  │   ├── category_api_test.go  # Category endpoint tests
  │   ├── transaction_api_test.go # Transaction endpoint tests
  │   └── sync_api_test.go      # Sync endpoint tests
  ├── integration/        # Integration tests
  │   └── sync_flow_integration_test.go # Complete sync flow tests
  ├── test_suite_test.go  # Test suite runner
  └── doc.go              # This documentation

Running Tests:

  # Run all tests
  go test -v ./test/...

  # Run unit tests only
  go test -v ./test/unit

  # Run API tests only
  go test -v ./test/api

  # Run integration tests only
  go test -v ./test/integration

  # Run with coverage
  go test -v -coverprofile=coverage.out ./test/...
  go tool cover -html=coverage.out

Test Categories:

1. Unit Tests
   - LWW (Last-Write-Wins) conflict resolution strategy
   - Edge cases: nil inputs, same timestamps, etc.

2. API Tests
   - Authentication (register/login)
   - Account CRUD operations
   - Category CRUD operations
   - Transaction CRUD operations
   - Sync pull/push endpoints
   - All tests use mocked dependencies

3. Integration Tests
   - Complete transaction creation and sync flow
   - Multi-device sync scenarios
   - Conflict resolution in real-world scenarios
   - WebSocket notification testing

See the server/README.md file for more detailed testing instructions.
*/
package test
