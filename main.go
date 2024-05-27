package main

import (
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

func HandleProcessData(ctx handle.Context[ProcessDataInput, ProcessDataParams]) (handle.User, error) {
	var user handle.User
	user.Name = ctx.Body.Name

	user.Age += ctx.Params.Id

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

	handle.Handle(r.Post, "/user/{id}", HandleProcessData, handle.NewHTMLTemplateWriter(tmpl, "test"))

	err = http.ListenAndServe(":8078", r)
	if err != nil {
		panic(err)
	}
}
