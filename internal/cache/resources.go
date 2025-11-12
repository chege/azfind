package cache

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Resource represents a cached Azure resource.
type Resource struct {
	ID             string
	Name           string
	Type           string
	SubscriptionID string
	ResourceGroup  string
	Location       string
	TenantID       string
	UpdatedAt      time.Time
}

// InsertResources inserts or replaces multiple resources transactionally.
func (db *DB) InsertResources(ctx context.Context, resources []Resource) error {
	if len(resources) == 0 {
		return nil
	}

	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT OR REPLACE INTO resources
		(id, name, type, subscriptionId, resourceGroup, location, tenantId, updatedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP);
	`)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("prepare insert: %w (rollback failed: %v)", err, rbErr)
		}
		return fmt.Errorf("prepare insert: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	for _, r := range resources {
		if _, err := stmt.ExecContext(ctx, r.ID, r.Name, r.Type, r.SubscriptionID, r.ResourceGroup, r.Location, r.TenantID); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("insert resource %q: %w (rollback failed: %v)", r.ID, err, rbErr)
			}
			return fmt.Errorf("insert resource %q: %w", r.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

// ListResources returns all resources ordered by name, resourceGroup, and type (all COLLATE NOCASE ASC).
func (db *DB) ListResources(ctx context.Context) ([]Resource, error) {
	query := `
		SELECT id, name, type, subscriptionId, resourceGroup, location, tenantId, updatedAt
		FROM resources
		ORDER BY name COLLATE NOCASE ASC,
		         resourceGroup COLLATE NOCASE ASC,
		         type COLLATE NOCASE ASC;`
	rows, err := db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query resources: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	return scanResources(rows)
}

// FindResources performs a LIKE search on name or id without limit.
func (db *DB) FindResources(ctx context.Context, query string) ([]Resource, error) {
	pattern := "%" + query + "%"
	baseQuery := `
		SELECT id, name, type, subscriptionId, resourceGroup, location, tenantId, updatedAt
		FROM resources
		WHERE name LIKE ? OR id LIKE ?
		ORDER BY name COLLATE NOCASE ASC,
		         resourceGroup COLLATE NOCASE ASC,
		         type COLLATE NOCASE ASC;`
	rows, err := db.conn.QueryContext(ctx, baseQuery, pattern, pattern)
	if err != nil {
		return nil, fmt.Errorf("query pattern: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	return scanResources(rows)
}

func (db *DB) FindResourcesByNamePrefix(ctx context.Context, name string) ([]Resource, error) {
	pattern := name + "%"
	query := `
        SELECT id, name, type, subscriptionId, resourceGroup, location, tenantId, updatedAt
        FROM resources
        WHERE name LIKE ?
        ORDER BY name COLLATE NOCASE ASC,
                 resourceGroup COLLATE NOCASE ASC,
                 type COLLATE NOCASE ASC;`

	rows, err := db.conn.QueryContext(ctx, query, pattern)
	if err != nil {
		return nil, fmt.Errorf("query by name: %w", err)
	}
	defer func() { _ = rows.Close() }()

	return scanResources(rows)
}

func (db *DB) FindResourceByExactName(ctx context.Context, name string) (*Resource, error) {
	query := `
        SELECT id, name, type, subscriptionId, resourceGroup, location, tenantId, updatedAt
        FROM resources
        WHERE LOWER(name) = LOWER(?)
        LIMIT 1;`

	row := db.conn.QueryRowContext(ctx, query, name)
	var r Resource
	if err := row.Scan(&r.ID, &r.Name, &r.Type, &r.SubscriptionID, &r.ResourceGroup, &r.Location, &r.TenantID, &r.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan resource: %w", err)
	}
	return &r, nil
}

func scanResources(rows *sql.Rows) ([]Resource, error) {
	var results []Resource
	for rows.Next() {
		var r Resource
		if err := rows.Scan(&r.ID, &r.Name, &r.Type, &r.SubscriptionID, &r.ResourceGroup, &r.Location, &r.TenantID, &r.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan resource: %w", err)
		}
		results = append(results, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iteration: %w", err)
	}
	return results, nil
}
