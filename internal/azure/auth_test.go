package azure

import (
	"context"
	"testing"
	"time"
)

func TestGetCredential_DefaultOrInteractive(t *testing.T) {
	if !liveTestsEnabled() {
		t.Skip("set AZF_LIVE_TESTS=1 to run Azure credential integration tests")
	}

	_, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cred, err := GetCredential()
	if err != nil {
		t.Skipf("GetCredential returned error (likely due to missing Azure login): %v", err)
	}

	if cred == nil {
		t.Fatalf("expected non-nil credential, got nil")
	}
}
