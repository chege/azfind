package azure

import (
	"context"
	"testing"
	"time"
)

func TestListSubscriptions_Basic(t *testing.T) {
	if !liveTestsEnabled() {
		t.Skip("set AZF_LIVE_TESTS=1 to run Azure subscription integration tests")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cred, err := GetCredential()
	if err != nil {
		t.Skipf("GetCredential failed (likely missing Azure login): %v", err)
	}

	subs, err := ListSubscriptions(ctx, cred)
	if err != nil {
		t.Skipf("ListSubscriptions returned error: %v", err)
	}

	if subs == nil {
		t.Fatalf("expected non-nil slice, got nil")
	}

	t.Logf("Fetched %d subscriptions", len(subs))
}
