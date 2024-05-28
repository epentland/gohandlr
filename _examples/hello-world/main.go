package main

import (
	"context"
	"net/http"
	"text/template"

	handle "github.com/epentland/twirp"
	options "github.com/epentland/twirp/options"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type ProcessDataInput struct {
	Name string `json:"name"`
}

type ProcessDataParams struct {
	Id  int `path:"id"`
	Age int `query:"age"`
}

func Process(ctx context.Context, body ProcessDataInput, params ProcessDataParams) (User, error) {
	// Do some processing
	user := User{
		Name: body.Name,
		Age:  params.Age,
	}
	return user, nil
}

func main() {
	// Create a text template
	tmplString := "<html><body>Hello, {{.Name}}, you are {{.Age}} years old!</body></html>"
	tmpl, err := template.New("index.html").Parse(tmplString)
	if err != nil {
		panic(err)
	}

	// Create a new http.ServeMux
	mux := http.NewServeMux()

	handle.Handle(mux.HandleFunc, "POST /user/{id}", Process,
		options.DefaultOptions(),
		options.WithJsonWriter(),
		options.WithHTMLTemplateWriter(tmpl, "index.html"),
	)

	err = http.ListenAndServe(":8082", mux)
	if err != nil {
		panic(err)
	}
}
