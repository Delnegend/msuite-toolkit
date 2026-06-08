package main

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"msuite-toolkit/pkg/app"
	add_users_to_group "msuite-toolkit/pkg/endpoints/add-users-to-group"
	finduserbyemail "msuite-toolkit/pkg/endpoints/find-user-by-email"
	"msuite-toolkit/pkg/types"
	"os"
	"sort"
	"sync"

	"github.com/alitto/pond/v2"
)

// resolvedUser pairs a configured email with the user ID it resolved to.
type resolvedUser struct {
	Email  string
	UserID string
}

func main() {
	app.Init("")

	as := &app.AppState
	requireGroupID(as)
	requireEmails(as)

	resolved, unresolved := resolveEmails(as)
	reportUnresolvedEmails(unresolved)

	if len(resolved) == 0 {
		slog.Error("no emails resolved to user IDs; nothing to add")
		os.Exit(1)
	}

	reportResolvedUsers(as, resolved)

	if as.DryRun {
		slog.Info("dry-run mode: users were NOT added to the group", "count", len(resolved))
		fmt.Printf("Would add %d user(s) to group %s — see report CSVs\n", len(resolved), as.GroupID)
		return
	}

	addUsers(as, resolved)
	fmt.Printf("Added %d user(s) to group %s — see report CSVs\n", len(resolved), as.GroupID)
}

// requireGroupID exits if the group ID is not set in config.
func requireGroupID(as *types.AppState) {
	if as.GroupID == "" {
		slog.Error("group_id is required in config")
		os.Exit(1)
	}
}

// requireEmails exits if no emails are configured.
func requireEmails(as *types.AppState) {
	if len(as.Emails) == 0 {
		slog.Error("emails is required in config")
		os.Exit(1)
	}
}

// resolveEmails resolves each configured email to a user ID in parallel and
// returns the successfully resolved users plus the emails that could not be
// resolved.
func resolveEmails(as *types.AppState) (resolved []resolvedUser, unresolved []string) {
	pool := pond.NewPool(as.WorkerCount)
	var mu sync.Mutex

	for _, email := range as.Emails {
		pool.Submit(func() {
			user, err := finduserbyemail.FindUserByEmail(as, email)
			if err != nil {
				slog.Error("could not resolve email to user ID", "email", email, "err", err)
				mu.Lock()
				unresolved = append(unresolved, email)
				mu.Unlock()
				return
			}
			mu.Lock()
			resolved = append(resolved, resolvedUser{Email: email, UserID: user.UserID})
			mu.Unlock()
			slog.Info("resolved email to user ID", "email", email, "user_id", user.UserID)
		})
	}

	pool.StopAndWait()

	sort.Slice(resolved, func(i, j int) bool { return resolved[i].Email < resolved[j].Email })
	sort.Strings(unresolved)
	return resolved, unresolved
}

// extractUserIDs pulls the UserID from each resolved user into a string slice.
func extractUserIDs(resolved []resolvedUser) []string {
	userIDs := make([]string, 0, len(resolved))
	for _, r := range resolved {
		userIDs = append(userIDs, r.UserID)
	}
	return userIDs
}

// addUsers performs the batched add and exits on failure.
func addUsers(as *types.AppState, resolved []resolvedUser) {
	userIDs := extractUserIDs(resolved)
	slog.Info("adding users to group", "group_id", as.GroupID, "count", len(userIDs))

	if err := add_users_to_group.AddUsersToGroupWithProgress(as, as.GroupID, userIDs); err != nil {
		slog.Error("adding users to group failed", "err", err)
		os.Exit(1)
	}
}

// reportResolvedUsers writes a CSV with the users that are candidates to be
// added (or were added). The file is named based on dry_run mode.
func reportResolvedUsers(as *types.AppState, resolved []resolvedUser) {
	csvFileName := "added-users.csv"
	if as.DryRun {
		csvFileName = "to-be-added-users.csv"
	}

	csvFile, err := os.Create(csvFileName)
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

	if err := w.Write([]string{"Email", "UserID", "GroupID"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}

	for _, r := range resolved {
		if err := w.Write([]string{r.Email, r.UserID, as.GroupID}); err != nil {
			slog.Error("writing csv row failed", "err", err)
			os.Exit(1)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		slog.Error("csv flush failed", "err", err)
		os.Exit(1)
	}

	slog.Info("wrote resolved users report", "csv", csvFileName, "count", len(resolved))
}

// reportUnresolvedEmails writes a CSV with the emails that could not be resolved
// to a user ID.
func reportUnresolvedEmails(unresolved []string) {
	if len(unresolved) == 0 {
		return
	}

	csvFileName := "unresolved-emails.csv"

	csvFile, err := os.Create(csvFileName)
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

	if err := w.Write([]string{"Email"}); err != nil {
		slog.Error("writing csv header failed", "err", err)
		os.Exit(1)
	}

	for _, email := range unresolved {
		if err := w.Write([]string{email}); err != nil {
			slog.Error("writing csv row failed", "err", err)
			os.Exit(1)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		slog.Error("csv flush failed", "err", err)
		os.Exit(1)
	}

	slog.Warn("wrote unresolved emails report", "csv", csvFileName, "count", len(unresolved))
}
