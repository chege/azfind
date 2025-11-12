package completion

import (
	"context"
	"strings"

	"github.com/chege/azfind/internal/cache"
)

// Generate prints resource name completions.
// 'partial' is what shell sends, but you may ignore it for now.
func Generate(ctx context.Context, partial string) ([]string, error) {
	db, err := cache.Open(ctx)
	if err != nil {
		// For completion handlers: fail quietly
		return nil, nil
	}
	defer func(db *cache.DB) {
		_ = db.Close()
	}(db)

	var resources []cache.Resource
	if partial == "" {
		resources, err = db.ListResources(ctx)
	} else {
		// partial match: allow substring, case-insensitive
		resources, err = db.FindResourcesByNamePrefix(ctx, partial)
	}
	if err != nil {
		return nil, nil // silent failure
	}

	seen := make(map[string]struct{}, len(resources))
	var results []string

	for _, r := range resources {
		// completion candidates MUST be unique
		if _, ok := seen[r.Name]; ok {
			continue
		}
		seen[r.Name] = struct{}{}

		// avoid printing empty names (rare)
		if strings.TrimSpace(r.Name) != "" {
			results = append(results, r.Name)
		}
	}

	return results, nil
}
