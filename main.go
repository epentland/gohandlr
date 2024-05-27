package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/epentland/twirp/handle"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ProcessDataInput struct {
	Name  string
	Email string
	Age   int
}

func ProcessData(ctx context.Context, params ProcessDataParams, input ProcessDataInput) (handle.User, error) {
	// Retrieve the user struct from the context
	user, ok := ctx.Value(handle.User{}).(handle.User)
	if !ok {
		return handle.User{}, fmt.Errorf("user not found in context")
	}

	// Do some processing
	user.Age += params.Id
	return user, nil
}

type ProcessDataParams struct {
	Id    int `path:"id"`
	Index int `query:"index"`
}

func main() {
	tmpl, err := template.ParseGlob("./templates/**/*.html")
	if err != nil {
		log.Fatal(err)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	handle.Handle(r.Post, "/user/{id}", ProcessData, handle.NewHTMLTemplateWriter(tmpl, "test"))

	err = http.ListenAndServe(":8087", r)
	if err != nil {
		panic(err)
	}
}
