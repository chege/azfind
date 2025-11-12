package cache

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestOpenAndCloseDB(t *testing.T) {
	ctx := context.Background()
	tmp := t.TempDir()
	_ = os.Setenv("XDG_CACHE_HOME", tmp)
	defer func() {
		_ = os.Unsetenv("XDG_CACHE_HOME")
	}()

	db, err := Open(ctx)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer func(db *DB) {
		_ = db.Close()
	}(db)

	// Check schema exists
	_, err = db.conn.ExecContext(ctx, "SELECT 1 FROM resources LIMIT 1;")
	if err != nil {
		t.Fatalf("expected resources table to exist: %v", err)
	}
}

func TestInsertAndQueryResources(t *testing.T) {
	ctx := context.Background()
	tmp := t.TempDir()
	_ = os.Setenv("XDG_CACHE_HOME", tmp)
	defer func() {
		_ = os.Unsetenv("XDG_CACHE_HOME")
	}()

	db, err := Open(ctx)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer func(db *DB) {
		_ = db.Close()
	}(db)

	resources := []Resource{
		{
			ID:             "1",
			Name:           "vm-prod",
			Type:           "Microsoft.Compute/virtualMachines",
			SubscriptionID: "sub1",
			ResourceGroup:  "rg1",
			Location:       "westeurope",
			TenantID:       "tenant1",
		},
		{
			ID:             "2",
			Name:           "storage-backup",
			Type:           "Microsoft.Storage/storageAccounts",
			SubscriptionID: "sub1",
			ResourceGroup:  "rg2",
			Location:       "norwayeast",
			TenantID:       "tenant1",
		},
	}

	if err := db.InsertResources(ctx, resources); err != nil {
		t.Fatalf("failed to insert resources: %v", err)
	}

	list, err := db.ListResources(ctx)
	if err != nil {
		t.Fatalf("failed to list resources: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(list))
	}

	found, err := db.FindResources(ctx, "vm")
	if err != nil {
		t.Fatalf("failed to find resources: %v", err)
	}
	if len(found) != 1 || found[0].Name != "vm-prod" {
		t.Fatalf("expected to find 'vm-prod', got %+v", found)
	}
}

func TestDBCreatesUnderCacheDir(t *testing.T) {
	ctx := context.Background()
	tmp := t.TempDir()
	_ = os.Setenv("XDG_CACHE_HOME", tmp)
	defer func() {
		_ = os.Unsetenv("XDG_CACHE_HOME")
	}()

	db, err := Open(ctx)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer func(db *DB) {
		_ = db.Close()
	}(db)

	dbPath := filepath.Join(tmp, "azf", "azf.db")
	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("expected db file at %s: %v", dbPath, err)
	}
}
