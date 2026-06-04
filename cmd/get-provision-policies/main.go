package main

import (
	"encoding/csv"
	"encoding/json"
	"log/slog"
	"msuite-toolkit/pkg/app"
	get_provision_policies "msuite-toolkit/pkg/endpoints/get-provision-policies"
	"msuite-toolkit/pkg/types"
	"os"
	"strconv"
)

func main() {
	outputPath := app.Init("provision_policies.csv")

	policies := fetchProvisionPolicies()
	writeProvisionPoliciesCSV(outputPath, policies)
}

// fetchProvisionPolicies retrieves all provision policies with progress.
func fetchProvisionPolicies() []types.ProvisionPolicy {
	policies, err := get_provision_policies.GetProvisionPoliciesWithProgress(&app.AppState)
	if err != nil {
		slog.Error("fetching provision policies failed", "err", err)
		os.Exit(1)
	}
	return policies
}

// marshalPolicyJSON serializes a policy to JSON, exiting on failure.
func marshalPolicyJSON(policy types.ProvisionPolicy) string {
	rawJSON, err := json.Marshal(policy)
	if err != nil {
		slog.Error("marshalling policy failed", "err", err, "policyID", policy.ProvisionPolicyID)
		os.Exit(1)
	}
	return string(rawJSON)
}

// writeProvisionPoliciesCSV writes the policies to a pipe-delimited CSV.
func writeProvisionPoliciesCSV(outputPath *string, policies []types.ProvisionPolicy) {
	csvFile, err := os.Create(*outputPath)
	if err != nil {
		slog.Error("creating csv file failed", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := csvFile.Close(); err != nil {
			slog.Error("closing csv file failed", "err", err)
		}
	}()

	w := csv.NewWriter(csvFile)
	w.Comma = '|'
	defer w.Flush()

	if err := w.Write([]string{"ProvisionPolicyID", "Name", "CreatedTime", "RawJSON"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}

	for _, policy := range policies {
		if err := w.Write([]string{
			policy.ProvisionPolicyID,
			policy.Name,
			strconv.FormatInt(policy.CreatedTime, 10),
			marshalPolicyJSON(policy),
		}); err != nil {
			slog.Error("writing csv row failed", "err", err)
			os.Exit(1)
		}
	}

	if err := w.Error(); err != nil {
		slog.Error("csv writer encountered error", "err", err)
		os.Exit(1)
	}
}
