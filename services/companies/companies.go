package companies

import "fmt"

// GetUser retrieves the user
func GetUser(users []Users, user string) (int, error) {
	for i, u := range users {
		if u.User == user {
			return i, nil
		}
	}
	return -1, fmt.Errorf("user not found")
}
