package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/epentland/gohandlr/examples/hello_world/handlr"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func HandlePostUser(ctx context.Context, req handlr.PostUsersIdInput) (handlr.User, error) {
	return handlr.User{
		Id:   req.Id,
		Name: req.Body.Name,
	}, nil
}

func main() {
	r := chi.NewMux()
	r.Use(middleware.Logger)

	r.MethodFunc(handlr.HandlePostUsersId(HandlePostUser))

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
	}
}
