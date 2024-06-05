package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type PostFilesImp struct{}

var _ PostFilesManager = PostFilesImp{}

func (i PostFilesImp) Handle(ctx context.Context, user User) (User, error) {
	user2 := User{
		Age: 15,
	}
	return user2, nil
}

func (i PostFilesImp) ReadApplicationJson(r *http.Request) (User, error) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	return user, nil
}

func (i PostFilesImp) Write200ApplicationJson(w http.ResponseWriter, user User) error {
	resp, err := json.Marshal(user)
	if err != nil {
		return err
	}
	w.Write(resp)
	return nil
}

func (i PostFilesImp) Write200TextHtml(w http.ResponseWriter, user User) error {
	resp, err := json.Marshal(user)
	if err != nil {
		return err
	}
	w.Write(resp)
	return nil
}

func main() {
	r := chi.NewMux()
	r.Use(middleware.Logger)
	r.MethodFunc(RegisterPostFiles(PostFilesImp{}))

	http.ListenAndServe(":8084", r)
}
