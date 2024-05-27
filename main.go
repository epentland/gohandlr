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
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func HandleProcessData(ctx context.Context, input ProcessDataInput, params ProcessDataParams) (handle.User, error) {
	fmt.Println(params, input)
	var user handle.User
	user.Name = input.Name
	user.Email = input.Email
	user.Age = input.Age

	user.Age += params.Id

	// Do some processing
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

	handle.Handle(r.Post, "/user/{id}", HandleProcessData, handle.NewHTMLTemplateWriter(tmpl, "test"), handle.NewJsonWriter())

	err = http.ListenAndServe(":8078", r)
	if err != nil {
		panic(err)
	}
}
