package syncer

import (
	"context"
	"fmt"
	"log"

	"github.com/chege/azfind/internal/azure"
	"github.com/chege/azfind/internal/cache"
)

func SyncAll(ctx context.Context) error {
	// Step 1: Authenticate
	cred, err := azure.GetCredential()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Step 2: Initialize cache
	db, err := cache.Open(ctx)
	if err != nil {
		return fmt.Errorf("failed to open cache: %w", err)
	}

	// Step 3: List subscriptions
	subs, err := azure.ListSubscriptions(ctx, cred)
	if err != nil {
		return fmt.Errorf("failed to list subscriptions: %w", err)
	}
	if len(subs) == 0 {
		fmt.Println("No subscriptions found.")
		return nil
	}

	total := 0
	for _, sub := range subs {
		if sub == nil || sub.SubscriptionID == nil {
			continue
		}

		subID := *sub.SubscriptionID
		fmt.Printf("Syncing subscription: %s\n", subID)

		resList, err := azure.ListResources(ctx, cred, subID, 5000)
		if err != nil {
			log.Printf("warning: failed to list resources for %s: %v\n", subID, err)
			continue
		}

		resources := make([]cache.Resource, 0, len(resList))
		for _, r := range resList {
			resources = append(resources, cache.Resource{
				ID:             fmt.Sprintf("%v", r["id"]),
				Name:           fmt.Sprintf("%v", r["name"]),
				Type:           fmt.Sprintf("%v", r["type"]),
				SubscriptionID: fmt.Sprintf("%v", r["subscriptionId"]),
				ResourceGroup:  fmt.Sprintf("%v", r["resourceGroup"]),
				Location:       fmt.Sprintf("%v", r["location"]),
				TenantID:       fmt.Sprintf("%v", r["tenantId"]),
			})
		}

		if err := db.InsertResources(ctx, resources); err != nil {
			return fmt.Errorf("sync: failed to insert resources for subscription %s: %w", subID, err)
		}

		total += len(resources)
		fmt.Printf("  → Synced %d resources\n", len(resources))
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close cache db: %w", err)
	}

	fmt.Printf("Sync completed. Total resources cached: %d\n", total)
	return nil
}
