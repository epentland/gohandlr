# gohandlr

`gohandlr` is a simple and intuitive to streamline the creation of HTTP handlers. It allows developers to focus on the essential data passing in and out of the handler, while providing standard interfaces for common tasks such as parsing request bodies, parameters, and response writers.

## Features

- Simplified HTTP handler creation
- Automatic parsing of request bodies and parameters
- Extensible options for customizing request handling
- Support for multiple content types and response formats
- Easy integration with existing Go HTTP libraries

## Installation

To install `gohandlr`, use the following command:

```shell
go get github.com/epentland/gohandlr
```

## Usage

Here's a simple example of how to use `gohandlr` to create an HTTP handler:

```go
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

// If the body, params or return struct are not needed, use the handle.Nil type.
func HandleNoBody(ctx context.Context, body gohandlr.Nil, params gohandlr.Nil) (gohandlr.Nil, error) {
	// Do write some data to the DB
	return gohandlr.Nil{}, nil
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

	gohandlr.Handle(mux.HandleFunc, "POST /user/{id}", HandleUserRequest,
		options.WithDefaults(),
		options.WithJsonWriter(),
		options.WithHTMLTemplateWriter(tmpl, "index.html"),
	)

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
```

In this example, we define a `User` struct to represent the response data, a `HandleUserBody` struct for the request body, and a `HandleUserParams` struct for the request parameters. The `HandleUserRequest` function is the actual HTTP handler, which takes the request body and parameters as input and returns a `User` instance.

To create the handler, we use the `gohandlr.Handle` function, passing in the `mux.HandleFunc` to register the handler, the HTTP method and path, the handler function, and any additional options (in this case, `options.WithDefaults()`).

If you don't need to use the request body, parameters, or return value, you can use the `gohandlr.Nil` type as a placeholder.

## Options

`gohandlr` provides a set of options to customize the behavior of the HTTP handler:

- `options.WithDefaults()`: Applies default options suitable for most use cases.
- `options.WithJsonWriter()`: Enables JSON response writing.
- `options.WithHTMLTemplateWriter(tmpl, name)`: Enables HTML template rendering for responses.
- `options.WithBodyReader()`: Specifies a custom request body reader.
- `options.WithParamsReader()`: Specifies a custom request parameter reader.

You can create your own options by implementing the appropriate interfaces:

- `BodyReader`: For parsing request bodies
- `ParamsReader`: For parsing request parameters
- `Writer`: For writing response data

## Contributing

Contributions to `gohandlr` are welcome! If you find a bug, have a feature request, or want to contribute code, please open an issue or submit a pull request on the [GitHub repository](https://github.com/epentland/gohandlr).

## License

`gohandlr` is released under the [MIT License](https://opensource.org/licenses/MIT).