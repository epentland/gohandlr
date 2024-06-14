package handlr

import (
	"context"
)

// POST request to /users/{id}
func processPostUsersId(ctx context.Context, req PostUsersIdInput) (User, error) {
	var resp User
	return resp, nil
}
