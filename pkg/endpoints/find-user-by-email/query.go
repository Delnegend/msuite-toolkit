package finduserbyemail

import (
	"fmt"
	get_users "msuite-toolkit/pkg/endpoints/get-users"
	"msuite-toolkit/pkg/types"
)

func FindUserByEmail(as *types.AppState, email string) (*types.UserInfo, error) {
	_, users, err := get_users.GetUsers(as, types.NewQueryRequestBuilder().WithSearch(email).WithFilters([]any{
		map[string]string{
			"key":   "email",
			"value": email,
		},
	}).Build())
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("User with email %s not found", email)
	}
	return &users[0], nil
}
