package handle

import (
	"context"
	"fmt"
	"net/http"
)

type Nil struct{}

type Options interface{}

type HandlerFunc[T, S any] func(context.Context, T) (S, error)

type HttpHandler func(string, http.HandlerFunc)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type HandleStruct[B, P, R any] struct {
	writers      map[string]Writer
	bodyReaders  map[string]BodyReader
	paramsReader ParamsReader
}

func NewHandleStruct[Body, Params, Response any](options ...Options) (*HandleStruct[Body, Params, Response], error) {
	data := &HandleStruct[Body, Params, Response]{
		writers:     make(map[string]Writer),
		bodyReaders: make(map[string]BodyReader),
	}

	// Add the default params reader
	data.paramsReader = WithParamsReader()

	// Add the default body reader
	data.bodyReaders["application/json"] = WithJSONBodyReader()

	for _, option := range options {
		switch v := option.(type) {
		case Writer:
			data.writers[v.Accept()] = v
		case BodyReader:
			data.bodyReaders[v.ContentType()] = v
		case ParamsReader:
			data.paramsReader = v
		default:
			return nil, fmt.Errorf("unknown option type %T", v)
		}
	}

	// Validation

	// Make sure there is a body reader if a body was passed
	var body Body
	if notNil(body) && len(data.bodyReaders) == 0 {
		return nil, fmt.Errorf("no body reader provided, please provide one")
	}

	// Make sure there is a params reader if a body was passed
	var params Params
	if notNil(params) && data.paramsReader == nil {
		return nil, fmt.Errorf("no params reader provided, please provide one")
	}

	// Make sure there is a writer if a response was passed
	var response Response
	if notNil(response) && len(data.writers) == 0 {
		return nil, fmt.Errorf("no writers provided, please provide one")

	}

	return data, nil
}

func notNil[T any](t T) bool {
	_, ok := any(t).(Nil)
	return !ok
}

func NewStruct[Body, Params, Response any](options ...Options) (*HandleStruct[Body, Params, Response], error) {
	data := &HandleStruct[Body, Params, Response]{
		writers:     make(map[string]Writer),
		bodyReaders: make(map[string]BodyReader),
	}

	// Add the default params reader
	data.paramsReader = WithParamsReader()

	// Add the default body reader
	data.bodyReaders["application/json"] = WithJSONBodyReader()

	for _, option := range options {
		switch v := option.(type) {
		case Writer:
			data.writers[v.Accept()] = v
		case BodyReader:
			data.bodyReaders[v.ContentType()] = v
		case ParamsReader:
			data.paramsReader = v
		default:
			return nil, fmt.Errorf("unknown option type %T", v)
		}
	}

	// Validation

	// Make sure there is a body reader if a body was passed
	var body Body
	if notNil(body) && len(data.bodyReaders) == 0 {
		return nil, fmt.Errorf("no body reader provided, please provide one")
	}

	// Make sure there is a params reader if a body was passed
	var params Params
	if notNil(params) && data.paramsReader == nil {
		return nil, fmt.Errorf("no params reader provided, please provide one")
	}

	// Make sure there is a writer if a response was passed
	var response Response
	if notNil(response) && len(data.writers) == 0 {
		return nil, fmt.Errorf("no writers provided, please provide one")

	}

	return data, nil
}

type ResponseFunc[B, P, R any] func(context.Context, B, P) (R, error)

func Handle[B, P, R any](httpHandler HttpHandler, path string, response ResponseFunc[B, P, R], options ...Options) {
	h, err := NewStruct[B, P, R](options...)
	if err != nil {
		panic(err)
	}

	httpHandler(path, func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		var body B
		var params P
		var err error

		// Read the request body if there are body readers
		if notNil(body) {
			if contentType == "" {
				contentType = "application/json"
			}

			// Check to see if the content type is in the bodyreader map
			bodyReader, ok := h.bodyReaders[contentType]
			if !ok {
				http.Error(w, fmt.Sprintf("unsupported content type %s", contentType), http.StatusUnsupportedMediaType)
				return
			}

			// Read the request body
			err = bodyReader.Reader(r, &body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		// Read the request params if there are params readers
		if notNil(params) {
			err = h.paramsReader.Reader(r, &params)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		// Does the data processes on the users function
		resp, err := response(r.Context(), body, params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check to see if there are any writers
		if notNil(resp) {
			// Get the Accept header
			accept := r.Header.Get("Accept")
			if accept == "" || accept == "*/*" {
				accept = "application/json"
			}

			// Check to see if the accept type is in the writers map
			writer, ok := h.writers[accept]
			if !ok {
				http.Error(w, fmt.Sprintf("unsupported accept type %s", accept), http.StatusNotAcceptable)
				return
			}

			// Write the response
			err = writer.Write(w, r, resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})
}

func (h HandleStruct[B, P, R]) Handle(httpHandler HttpHandler, path string, response ResponseFunc[B, P, R]) {
	httpHandler(path, func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		var body B
		var params P
		var err error

		// Read the request body if there are body readers
		if notNil(body) {
			if contentType == "" {
				contentType = "application/json"
			}

			// Check to see if the content type is in the bodyreader map
			bodyReader, ok := h.bodyReaders[contentType]
			if !ok {
				http.Error(w, fmt.Sprintf("unsupported content type %s", contentType), http.StatusUnsupportedMediaType)
				return
			}

			// Read the request body
			err = bodyReader.Reader(r, &body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		// Read the request params if there are params readers
		if notNil(params) {
			err = h.paramsReader.Reader(r, &params)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		// Does the data processes on the users function
		resp, err := response(r.Context(), body, params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check to see if there are any writers
		if notNil(resp) {
			// Get the Accept header
			accept := r.Header.Get("Accept")
			if accept == "" || accept == "*/*" {
				accept = "application/json"
			}

			// Check to see if the accept type is in the writers map
			writer, ok := h.writers[accept]
			if !ok {
				http.Error(w, fmt.Sprintf("unsupported accept type %s", accept), http.StatusNotAcceptable)
				return
			}

			// Write the response
			err = writer.Write(w, r, resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})

}
