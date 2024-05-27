package handle

import (
	"context"
	"net/http"
)

type HandlerFunc[T, S any] func(context.Context, T) (S, error)

type ProcessFunc[B, P, R any] func(Context[B, P]) (R, error)

type HttpHandler func(string, http.HandlerFunc)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type Context[Body, Params any] struct {
	Context context.Context
	Body    Body
	Params  Params
}

func Handle[B, P, S any](httpHandler HttpHandler, path string, process ProcessFunc[B, P, S], acceptWriters ...Writer) {
	httpHandler(path, func(w http.ResponseWriter, r *http.Request) {

		// Decode the Request params
		params, err := DecodeRequestParams[P](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Decode the Request input
		body, err := DecodeRequestBody[B](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create the context
		ctx := Context[B, P]{
			Context: r.Context(),
			Body:    body,
			Params:  params,
		}

		// Process the data
		resp, err := process(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Checks the Accept header
		accept := r.Header.Get("Accept")

		if accept == "" || accept == "*/*" {
			accept = "application/json"
		}

		// Write the response
		for _, writer := range acceptWriters {
			if writer.Accept() == accept {
				err = writer.Write(w, r, resp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}
		}

		// If no writer was found, return an error
		http.Error(w, "No writer found for Accept header: " + accept, http.StatusNotAcceptable)
	})
}
