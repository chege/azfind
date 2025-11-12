package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcegraph/armresourcegraph"
)

func ListResources(ctx context.Context, cred azcore.TokenCredential, subscriptionID string, limit int32) ([]map[string]any, error) {
	client, err := armresourcegraph.NewClient(cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource graph client: %w", err)
	}

	query := fmt.Sprintf("Resources | project id,name,type,subscriptionId,resourceGroup,location,tenantId | limit %d", limit)
	request := armresourcegraph.QueryRequest{
		Subscriptions: []*string{&subscriptionID},
		Query:         &query,
	}

	resp, err := client.Resources(ctx, request, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute resource graph query: %w", err)
	}

	if resp.Data == nil {
		return nil, fmt.Errorf("resource graph query returned no data")
	}

	data, ok := resp.Data.([]any)
	if !ok {
		return nil, fmt.Errorf("unexpected data format in resource graph response")
	}

	results := make([]map[string]any, 0, len(data))
	for _, item := range data {
		m, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("unexpected item format in resource graph response data")
		}
		results = append(results, m)
	}

	return results, nil
}
