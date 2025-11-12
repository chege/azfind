package azure

import (
	"context"
	"testing"
	"time"
)

func TestGetCredential_DefaultOrInteractive(t *testing.T) {
	_, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cred, err := GetCredential()
	if err != nil {
		t.Logf("GetCredential returned error (likely due to missing Azure login): %v", err)
	} else if cred == nil {
		t.Fatalf("expected non-nil credential, got nil")
	}
}
