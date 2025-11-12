package search

import (
	"context"
	"fmt"

	"github.com/chege/azfind/internal/cache"
	"github.com/rodaine/table"
)

// ListCache prints all cached resources in a concise, formatted table.
func ListCache(ctx context.Context) (err error) {
	db, err := cache.Open(ctx)
	if err != nil {
		return fmt.Errorf("open cache: %w", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close cache db: %w", cerr)
		}
	}()

	resources, err := db.FindResources(ctx, "")
	if err != nil {
		return fmt.Errorf("find resources: %w", err)
	}

	if len(resources) == 0 {
		fmt.Println("------------------------------------------------------------")
		fmt.Println("No cached resources found. Run 'azf sync' first.")
		fmt.Println("------------------------------------------------------------")
		return nil
	}

	fmt.Printf("Cached %d resources:\n\n", len(resources))

	tbl := table.New("Name", "Type", "Resource Group")
	for _, r := range resources {
		tbl.AddRow(r.Name, r.Type, r.ResourceGroup)
	}
	tbl.Print()

	return nil
}
