package handle

import (
	"context"
	"net/http"
)

type HandlerFunc[T, S any] func(context.Context, T) (S, error)

type HttpHandler func(string, http.HandlerFunc)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func Handle[T, P, S any](handle HttpHandler, path string, process func(context.Context, T, P) (S, error), acceptWriters ...Writer) {
	handle(path, func(w http.ResponseWriter, r *http.Request) {
		// Decode the Request input
		body, params, err := decodeRequest[T, P](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

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

		// Add json if not present
		hasJson := false
		for _, writer := range acceptWriters {
			if writer.Accept() == "application/json" {
				hasJson = true
				break
			}
		}

		if !hasJson {
			acceptWriters = append(acceptWriters, NewJsonWriter())
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
	})
}
