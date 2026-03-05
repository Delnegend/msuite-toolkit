package get_provision_policies

import (
	"msuite-toolkit/pkg/types"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestGetProvisionPolicies(t *testing.T) {
	var appState types.AppState
	if _, err := toml.DecodeFile("../../../config.test.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}

	resp, err := GetProvisionPolicies(&appState, 0, 100)
	if err != nil {
		t.Fatalf("GetProvisionPolicies failed: %v", err)
	}

	t.Logf("Provision policies count: %d (total: %d)", len(resp.Data), resp.Count)
	for i, policy := range resp.Data {
		if i < 5 {
			t.Logf("Policy %d: ID=%s, Name=%s", i, policy.ProvisionPolicyID, policy.Name)
		}
	}
}

func TestGetProvisionPoliciesWithProgress(t *testing.T) {
	var appState types.AppState
	if _, err := toml.DecodeFile("../../../config.test.toml", &appState); err != nil {
		t.Fatalf("decoding config file failed: %v", err)
	}

	policies, err := GetProvisionPoliciesWithProgress(&appState)
	if err != nil {
		t.Fatalf("GetProvisionPoliciesWithProgress failed: %v", err)
	}

	t.Logf("Total provision policies fetched with progress: %d", len(policies))
}
