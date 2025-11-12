package fzfui

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/chege/azfind/internal/cache"
)

// RunSearch performs optional prefiltering, launches fzf, and opens the selected resource.
func RunSearch(ctx context.Context, args []string) error {
	db, err := cache.Open(ctx)
	if err != nil {
		return fmt.Errorf("open cache: %w", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			fmt.Printf("warning: failed to close cache db: %v\n", cerr)
		}
	}()

	if len(args) == 1 {
		exactResource, err := db.FindResourceByExactName(ctx, args[0])
		if err != nil {
			return fmt.Errorf("find resource by exact name: %w", err)
		}
		if exactResource != nil {
			if err := openResource(*exactResource); err != nil {
				return fmt.Errorf("failed to open resource: %w", err)
			}
			return nil
		}
	}

	var resources []cache.Resource
	query := strings.Join(args, " ")

	if query == "" {
		// No filter → list all
		resources, err = db.ListResources(ctx)
		if err != nil {
			return fmt.Errorf("list resources: %w", err)
		}
		if len(resources) == 0 {
			fmt.Println("No cached resources found. Run `azf sync` first.")
			return nil
		}
	} else {
		// Load full list; only use query as initial filter for fzf
		resources, err = db.ListResources(ctx)
		if err != nil {
			return fmt.Errorf("list resources: %w", err)
		}
	}

	// If only one result remains → open directly
	if len(resources) == 1 {
		if err := openResource(resources[0]); err != nil {
			return fmt.Errorf("failed to open resource: %w", err)
		}
		return nil
	}

	// Run fzf selector
	selected, err := SelectResource(resources, query)
	if err != nil {
		return err
	}
	if selected == nil {
		// user cancelled
		return nil
	}

	if err := openResource(*selected); err != nil {
		return fmt.Errorf("failed to open resource: %w", err)
	}
	return nil
}

func openResource(r cache.Resource) error {
	url := fmt.Sprintf("https://portal.azure.com/#@%s/resource%s", r.TenantID, r.ID)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return fmt.Errorf("unsupported platform")
	}
	return cmd.Start()
}
