package get_provision_policies

import (
	"fmt"
	"msuite-toolkit/pkg/types"
)

// GetProvisionPoliciesWithProgress fetches all provision policies.
func GetProvisionPoliciesWithProgress(appState *types.AppState) ([]types.ProvisionPolicy, error) {
	fmt.Println("Fetching provision policies...")

	var allPolicies []types.ProvisionPolicy
	offset := 0
	limit := 200

	for {
		resp, err := GetProvisionPolicies(appState, offset, limit)
		if err != nil {
			return nil, fmt.Errorf("fetching provision policies failed: %w", err)
		}

		allPolicies = append(allPolicies, resp.Data...)

		if offset+len(resp.Data) >= resp.Count || len(resp.Data) == 0 {
			break
		}
		offset += limit
		fmt.Printf("\rFetched %d/%d policies...", len(allPolicies), resp.Count)
	}
	fmt.Printf("\rFetched %d policies.Done.\n", len(allPolicies))

	return allPolicies, nil
}
