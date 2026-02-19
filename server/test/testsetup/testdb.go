package testsetup

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// TestDBConfig holds test database configuration
type TestDBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// GetTestDBConfig returns test database configuration from env vars or defaults
func GetTestDBConfig() TestDBConfig {
	return TestDBConfig{
		Host:     getEnv("TEST_DB_HOST", "localhost"),
		Port:     getEnv("TEST_DB_PORT", "5432"),
		User:     getEnv("TEST_DB_USER", "postgres"),
		Password: getEnv("TEST_DB_PASSWORD", "postgres"),
		DBName:   getEnv("TEST_DB_NAME", "account_test"),
		SSLMode:  getEnv("TEST_DB_SSLMODE", "disable"),
	}
}

// ConnectionString returns the database connection string
func (c TestDBConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// SetupTestDB creates a test database connection and runs migrations
func SetupTestDB(t *testing.T) (*sqlx.DB, func()) {
	config := GetTestDBConfig()

	// Connect to postgres database to create test database
	postgresConfig := config
	postgresConfig.DBName = "postgres"

	db, err := sql.Open("postgres", postgresConfig.ConnectionString())
	if err != nil {
		t.Fatalf("Failed to connect to postgres: %v", err)
	}

	// Create test database
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", config.DBName))
	if err != nil {
		log.Printf("Warning: Failed to drop test database: %v", err)
	}
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", config.DBName))
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	db.Close()

	// Connect to test database
	testDB, err := sqlx.Connect("postgres", config.ConnectionString())
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	if err := runMigrations(testDB); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		testDB.Close()
		// Drop test database
		db, err := sql.Open("postgres", postgresConfig.ConnectionString())
		if err == nil {
			defer db.Close()
			_, _ = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", config.DBName))
		}
	}

	return testDB, cleanup
}

// runMigrations creates the database schema
func runMigrations(db *sqlx.DB) error {
	schema := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE users (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE TABLE accounts (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		name VARCHAR(255) NOT NULL,
		type VARCHAR(50) NOT NULL,
		currency VARCHAR(10) DEFAULT 'CNY',
		balance DECIMAL(15,2) DEFAULT 0,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		version INTEGER DEFAULT 1,
		is_deleted BOOLEAN DEFAULT FALSE
	);

	CREATE TABLE categories (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		name VARCHAR(255) NOT NULL,
		type VARCHAR(50) NOT NULL,
		parent_id UUID REFERENCES categories(id) ON DELETE SET NULL,
		icon VARCHAR(50),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		version INTEGER DEFAULT 1,
		is_deleted BOOLEAN DEFAULT FALSE
	);

	CREATE TABLE transactions (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
		category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
		type VARCHAR(50) NOT NULL,
		amount DECIMAL(15,2) NOT NULL,
		currency VARCHAR(10) DEFAULT 'CNY',
		note TEXT,
		transaction_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		version INTEGER DEFAULT 1,
		is_deleted BOOLEAN DEFAULT FALSE
	);

	CREATE TABLE transfer_links (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		from_transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
		to_transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE TABLE sync_state (
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		device_id VARCHAR(255) NOT NULL,
		last_sync_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		sync_token VARCHAR(255),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		PRIMARY KEY (user_id, device_id)
	);

	CREATE INDEX idx_accounts_user_id ON accounts(user_id);
	CREATE INDEX idx_accounts_last_modified ON accounts(user_id, last_modified_at);
	CREATE INDEX idx_categories_user_id ON categories(user_id);
	CREATE INDEX idx_categories_last_modified ON categories(user_id, last_modified_at);
	CREATE INDEX idx_transactions_user_id ON transactions(user_id);
	CREATE INDEX idx_transactions_account_id ON transactions(account_id);
	CREATE INDEX idx_transactions_date ON transactions(user_id, transaction_date);
	CREATE INDEX idx_transactions_last_modified ON transactions(user_id, last_modified_at);
	`

	_, err := db.Exec(schema)
	return err
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
