package handlr

import (
	"context"
	"github.com/go-chi/chi/v5"
)

func RegisterHandlers(r *chi.Mux) {
	r.MethodFunc(HandlePutUsersId())

}

// PUT request to /users/{id}
func processPutUsersId(ctx context.Context, req PutUsersIdInput) (User, error) {
	var resp User
	return resp, nil
}
