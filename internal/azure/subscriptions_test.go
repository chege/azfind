package azure

import (
	"context"
	"testing"
	"time"
)

func TestListSubscriptions_Basic(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cred, err := GetCredential()
	if err != nil {
		t.Logf("GetCredential failed (likely missing Azure login): %v", err)
		return
	}

	subs, err := ListSubscriptions(ctx, cred)
	if err != nil {
		t.Logf("ListSubscriptions returned error: %v", err)
		return
	}

	if subs == nil {
		t.Fatalf("expected non-nil slice, got nil")
	}

	t.Logf("Fetched %d subscriptions", len(subs))
}
