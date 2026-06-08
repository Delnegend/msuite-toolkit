package add_users_to_group

import (
	"fmt"
	"log/slog"
	"msuite-toolkit/pkg/types"
	"msuite-toolkit/pkg/utils"
	"sync"
)

// AddUsersToGroupWithProgress adds the provided users to the given group in
// batches while rendering a progress bar to stdout.
func AddUsersToGroupWithProgress(as *types.AppState, groupID string, userIDs []string) error {
	fmt.Println("Adding users to group...")
	var wg sync.WaitGroup
	progressPercentChan := make(chan int)
	wg.Go(func() {
		for percent := range progressPercentChan {
			utils.PrintProgressBar(percent)
		}
	})
	err := AddUsersToGroup(as, groupID, userIDs, progressPercentChan)
	close(progressPercentChan)
	wg.Wait()
	if err != nil {
		return fmt.Errorf("failed to add users to group: %w", err)
	}
	slog.Info("added users to group", "group_id", groupID, "count", len(userIDs))
	return nil
}
