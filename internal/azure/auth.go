package azure

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func GetCredential() (azcore.TokenCredential, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err == nil {
		fmt.Println("Authenticated using cached or default credentials.")
		return cred, nil
	}

	fmt.Println("Default credentials not available; opening browser for login...")
	interactive, ierr := azidentity.NewInteractiveBrowserCredential(nil)
	if ierr != nil {
		return nil, fmt.Errorf("failed to get Azure credentials: %w", ierr)
	}

	return interactive, nil
}
