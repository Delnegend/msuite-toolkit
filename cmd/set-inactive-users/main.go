package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"log/slog"
	"msuite-toolkit/pkg/app"
	inactive_user "msuite-toolkit/pkg/endpoints/inactive-user"
	"msuite-toolkit/pkg/types"
	"os"
	"strings"
)

func main() {
	inputPath := flag.String("input", "", "path to input file with user ids (one per line)")

	// call app.Init after registering our custom flag so it gets parsed together
	outputPath := app.Init("inactive_users.csv")

	if *inputPath == "" {
		slog.Error("input flag is required")
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(*inputPath)
	if err != nil {
		slog.Error("opening input file failed", "path", *inputPath, "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("closing input file failed", "err", err)
		}
	}()

	scanner := bufio.NewScanner(f)
	var total int
	var failed int
	// collect results for CSV output
	var rows [][]string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		total++
		if err := inactive_user.InactiveUser(&app.AppState, types.UserID(line)); err != nil {
			msg := err.Error()
			slog.Error("failed to inactive user", "user", line, "err", err)
			failed++
			rows = append(rows, []string{line, msg})
			continue
		}
		slog.Info("inactivated user", "user", line)
		rows = append(rows, []string{line, "OK"})
	}
	if err := scanner.Err(); err != nil {
		slog.Error("reading input file failed", "err", err)
		os.Exit(1)
	}

	// write results to output CSV (pipe-delimited)
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

	if err := w.Write([]string{"UserID", "Result"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}
	for _, r := range rows {
		if err := w.Write(r); err != nil {
			slog.Error("writing csv row failed", "err", err)
			os.Exit(1)
		}
	}

	if err := w.Error(); err != nil {
		slog.Error("csv writer encountered error", "err", err)
		os.Exit(1)
	}

	if failed > 0 {
		slog.Error("completed with failures", "processed", total, "failed", failed)
		os.Exit(1)
	}
	slog.Info("completed", "processed", total)
}
