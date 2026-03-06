package app

import (
	"flag"
	"log/slog"
	"os"

	"msuite-toolkit/pkg/types"

	"github.com/BurntSushi/toml"
)

var AppState types.AppState

// Init registers flags (including an `output` flag with the provided default),
// parses them, loads the config TOML into AppState and returns the output path pointer.
func Init(defaultOutput string) *string {
	configPath := flag.String("config", "./config.toml", "path to config file")
	outputPath := flag.String("output", defaultOutput, "path to output CSV file")
	helpFlag := flag.Bool("h", false, "show help")
	helpFlagLong := flag.Bool("help", false, "show help")

	flag.Parse()

	if *helpFlag || *helpFlagLong {
		flag.Usage()
		os.Exit(0)
	}

	if _, err := os.Stat(*configPath); err != nil {
		if os.IsNotExist(err) {
			slog.Error("config file not found", "path", *configPath, "err", err)
		} else {
			slog.Error("stat config file failed", "err", err)
		}
		os.Exit(1)
	}

	if _, err := toml.DecodeFile(*configPath, &AppState); err != nil {
		slog.Error("decoding config file failed", "err", err)
		os.Exit(1)
	}

	if AppState.WorkerCount == 0 {
		AppState.WorkerCount = 100
	}

	slog.Info("loaded config", "admin_portal", AppState.AdminPortalAddress, "admin_user_id", AppState.AdminUserID, "bearer_token_len", len(AppState.BearerToken), "worker_count", AppState.WorkerCount)

	return outputPath
}
