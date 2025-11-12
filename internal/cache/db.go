package cache

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// DB wraps the underlying SQL database connection.
type DB struct {
	conn *sql.DB
}

// Open initializes (or creates) the azf cache database.
func Open(ctx context.Context) (*DB, error) {
	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if cacheDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to resolve home dir: %w", err)
		}
		cacheDir = filepath.Join(home, ".cache")
	}

	dbPath := filepath.Join(cacheDir, "azf", "azf.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("failed to create cache dir: %w", err)
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open cache db: %w", err)
	}

	if err := conn.PingContext(ctx); err != nil {
		closeErr := conn.Close()
		if closeErr != nil {
			return nil, fmt.Errorf("failed to connect to cache db: %w (also failed to close: %v)", err, closeErr)
		}
		return nil, fmt.Errorf("failed to connect to cache db: %w", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS resources (
		id TEXT PRIMARY KEY,
		name TEXT,
		type TEXT,
		subscriptionId TEXT,
		resourceGroup TEXT,
		location TEXT,
		tenantId TEXT,
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := conn.ExecContext(ctx, schema); err != nil {
		closeErr := conn.Close()
		if closeErr != nil {
			return nil, fmt.Errorf("failed to initialize schema: %w (also failed to close: %v)", err, closeErr)
		}
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return &DB{conn: conn}, nil
}

// Close safely closes the database connection.
func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}
