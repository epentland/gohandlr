package main

import (
	"context"
	"net/http"
	"text/template"

	"github.com/epentland/gohandlr"
	"github.com/epentland/gohandlr/options"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type HandleUserBody struct {
	Name string `json:"name"`
}

type HandleUserParams struct {
	Id  int `path:"id"`
	Age int `query:"age"`
}

func HandleUserRequest(ctx context.Context, body HandleUserBody, params HandleUserParams) (User, error) {
	// Do some processing
	user := User{
		Name: body.Name,
		Age:  params.Age,
	}
	return user, nil
}

func HandleNoBody(ctx context.Context, body gohandlr.Nil, params gohandlr.Nil) (gohandlr.Nil, error) {
	// Do write some processing that doesn't require a return value
	return gohandlr.Nil{}, nil
}

func main() {
	// Create a html template
	tmplString := "<html><body>Hello, {{.Name}}, you are {{.Age}} years old!</body></html>"
	tmpl, err := template.New("index.html").Parse(tmplString)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

    // Works with any router
	gohandlr.Handle(mux.HandleFunc, "POST /user/{id}", HandleUserRequest,
		options.WithDefaults(),
		options.WithJsonWriter(),
		options.WithHTMLTemplateWriter(tmpl, "index.html"),
	)

    gohandlr.Handle(mux.HandleFunc, "PUT /user", HandleNoBody, options.WithDefaults())

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}