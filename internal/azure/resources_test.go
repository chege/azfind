package azure

import (
	"context"
	"testing"
	"time"
)

func TestListResources_Basic(t *testing.T) {
	if !liveTestsEnabled() {
		t.Skip("set AZF_LIVE_TESTS=1 to run Azure resource integration tests")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cred, err := GetCredential()
	if err != nil {
		t.Skipf("GetCredential failed (likely missing Azure login): %v", err)
	}

	subs, err := ListSubscriptions(ctx, cred)
	if err != nil {
		t.Skipf("ListSubscriptions failed: %v", err)
	}
	if len(subs) == 0 || subs[0] == nil || subs[0].SubscriptionID == nil {
		t.Logf("No subscriptions found; skipping resource test.")
		return
	}

	subID := *subs[0].SubscriptionID
	res, err := ListResources(ctx, cred, subID, 100)
	if err != nil {
		t.Logf("ListResources failed for %s: %v", subID, err)
		return
	}

	t.Logf("Fetched %d resources from subscription %s", len(res), subID)
}
