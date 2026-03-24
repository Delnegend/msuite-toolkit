package main

import (
	"encoding/csv"
	"encoding/json"
	"log/slog"
	"msuite-toolkit/pkg/app"
	get_provision_policies "msuite-toolkit/pkg/endpoints/get-provision-policies"
	"os"
	"strconv"
)

func main() {
	outputPath := app.Init("provision_policies.csv")

	as := &app.AppState

	policies, err := get_provision_policies.GetProvisionPoliciesWithProgress(as)
	if err != nil {
		slog.Error("fetching provision policies failed", "err", err)
		os.Exit(1)
	}

	// create CSV file
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

	// write header
	if err := w.Write([]string{"ProvisionPolicyID", "Name", "CreatedTime", "RawJSON"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}

	for _, policy := range policies {
		rawJSON, err := json.Marshal(policy)
		if err != nil {
			slog.Error("marshalling policy failed", "err", err, "policyID", policy.ProvisionPolicyID)
			os.Exit(1)
		}

		if err := w.Write([]string{
			policy.ProvisionPolicyID,
			policy.Name,
			strconv.FormatInt(policy.CreatedTime, 10),
			string(rawJSON),
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
