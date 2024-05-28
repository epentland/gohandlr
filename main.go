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

type ProcessData struct{}

func Process(ctx context.Context, body ProcessDataInput, params ProcessDataParams) (handle.Nil, error) {
	fmt.Println("Processing data")
	var user handle.User
	user.Name = body.Name
	user.Age = body.Age

	user.Age += 100 + params.Id
	return handle.Nil{}, nil
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
	fmt.Println(tmpl)

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.RequestID)

		handle.Handle(r.Post, "/user/{id}", Process,
			handle.WithJsonWriter(),
			handle.WithHTMLTemplateWriter(tmpl, "test"))
	})

	err = http.ListenAndServe(":8077", r)
	if err != nil {
		panic(err)
	}
}
