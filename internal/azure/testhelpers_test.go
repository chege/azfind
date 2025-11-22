package azure

import "os"

func liveTestsEnabled() bool {
	return os.Getenv("AZF_LIVE_TESTS") == "1"
}
