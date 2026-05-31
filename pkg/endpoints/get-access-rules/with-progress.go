package get_access_rules

import (
	"fmt"
	"msuite-toolkit/pkg/types"
)

// GetAccessRulesWithProgress fetches all access rules, paginating until complete.
func GetAccessRulesWithProgress(as *types.AppState) ([]types.AccessRule, error) {
	fmt.Println("Fetching access rules...")

	var all []types.AccessRule
	builder := types.NewQueryRequestBuilder()
	offset := 0
	limit := 200

	for {
		builder.WithOffset(offset).WithLimit(limit)
		payload := builder.Build()

		resp, err := GetAccessRules(as, payload)
		if err != nil {
			return nil, fmt.Errorf("fetching access rules failed: %w", err)
		}

		all = append(all, resp.Data...)

		if offset+len(resp.Data) >= resp.Count || len(resp.Data) == 0 {
			break
		}
		offset += limit
		fmt.Printf("\rFetched %d/%d access rules...", len(all), resp.Count)
	}
	fmt.Printf("\rFetched %d access rules.Done.\n", len(all))

	return all, nil
}
