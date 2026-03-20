package get_users

import (
	"fmt"
	"log/slog"
	"msuite-toolkit/pkg/types"
	"msuite-toolkit/pkg/utils"
	"os"
	"sync"
)

func GetUsersWithProgress(appState *types.AppState, basePayload types.GetUsersRequestPayload) []types.UserInfo {
	fmt.Println("Fetching users...")
	var wg sync.WaitGroup
	progressPercentChan := make(chan int)
	wg.Go(func() {
		for percent := range progressPercentChan {
			utils.PrintProgressBar(percent)
		}
	})
	users, err := GetAllUsers(appState, basePayload, progressPercentChan)
	close(progressPercentChan)
	wg.Wait()
	if err != nil {
		slog.Error("failed to get all users", "err", err)
		os.Exit(1)
	}
	slog.Info("fetched users", "count", len(users))
	return users
}
