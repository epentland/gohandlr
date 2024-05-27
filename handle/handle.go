package handle

import (
	"context"
	"fmt"
	"net/http"
)

type HandlerFunc[T, S any] func(context.Context, T) (S, error)

type ProcessFunc[T, P, S any] func(context.Context, T, P) (S, error)

type HttpHandler func(string, http.HandlerFunc)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func Handle[T, P, S any](handle HttpHandler, path string, process ProcessFunc[T, P, S], acceptWriters ...Writer) {
	handle(path, func(w http.ResponseWriter, r *http.Request) {
		var body T

		// Decode the Request input
		err := DecodeRequestBody(r, &body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(body)

		// Decode the Request params
		params, err := DecodeRequestParams[P](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(params)

		// Process the data
		resp, err := process(r.Context(), body, params)
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
		http.Error(w, "No writer found for Accept header", http.StatusNotAcceptable)
	})
}
