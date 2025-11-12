package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

// ListSubscriptions retrieves all accessible subscriptions using the provided credential.
func ListSubscriptions(ctx context.Context, cred azcore.TokenCredential) ([]*armsubscriptions.Subscription, error) {
	client, err := armsubscriptions.NewClient(cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscriptions client: %w", err)
	}

	pager := client.NewListPager(nil)
	var subs []*armsubscriptions.Subscription

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next subscription page: %w", err)
		}
		subs = append(subs, page.Value...)
	}

	return subs, nil
}
