package handle

import (
	"context"
	"fmt"
	"net/http"

	"github.com/epentland/twirp/options"
)

type Nil struct{}

type ResponseFunc[B, P, R any] func(context.Context, B, P) (R, error)

type Handler[Body, Params, Response any] struct {
	writers      map[string]options.Writer
	bodyReaders  map[string]options.BodyReader
	paramsReader options.ParamsReader
}

func notNil[T any](t T) bool {
	_, ok := interface{}(t).(Nil)
	return !ok
}

func ValidateOptions[Body, Params, Response any](opts ...options.Options) error {
	var body Body
	var params Params
	var response Response

	var hasBodyReader bool
	var hasParamsReader bool
	var hasWriter bool

	for _, opt := range opts {
		switch opt.(type) {
		case options.BodyReader:
			hasBodyReader = true
		case options.ParamsReader:
			hasParamsReader = true
		case options.Writer:
			hasWriter = true
		default:
			return fmt.Errorf("unknown option")
		}
	}

	if notNil(body) && !hasBodyReader {
		return fmt.Errorf("no body reader provided, please provide one")
	}

	if notNil(params) && !hasParamsReader {
		return fmt.Errorf("no params reader provided, please provide one")
	}

	if notNil(response) && !hasWriter {
		return fmt.Errorf("no writers provided, please provide one")
	}

	return nil
}

func SetOptions[Body, Params, Response any](data *Handler[Body, Params, Response], opts ...options.Options) {
	for _, opt := range opts {
		switch v := opt.(type) {
		case options.Writer:
			data.writers[v.Accept()] = v
		case options.BodyReader:
			data.bodyReaders[v.ContentType()] = v
		case options.ParamsReader:
			data.paramsReader = v
		}
	}
}

func NewHandler[Body, Params, Response any](opts ...options.Options) (*Handler[Body, Params, Response], error) {
	err := ValidateOptions[Body, Params, Response](opts...)
	if err != nil {
		return nil, err
	}

	handler := &Handler[Body, Params, Response]{
		writers:     make(map[string]options.Writer),
		bodyReaders: make(map[string]options.BodyReader),
	}

	SetOptions(handler, opts...)

	return handler, nil
}

func Handle[Body, Params, Response any](httpHandler func(string, func(http.ResponseWriter, *http.Request)), path string, handler ResponseFunc[Body, Params, Response], opts ...options.Options) error {
	h, err := NewHandler[Body, Params, Response](opts...)
	if err != nil {
		panic(err)
	}

	httpHandler(path, func(w http.ResponseWriter, r *http.Request) {
		handleRequest(h, w, r, handler)
	})

	return nil
}

func handleRequest[Body, Params, Response any](h *Handler[Body, Params, Response], w http.ResponseWriter, r *http.Request, handler ResponseFunc[Body, Params, Response]) {
	contentType := r.Header.Get("Content-Type")
	var body Body
	var params Params
	var err error

	if notNil(body) {
		body, err = readRequestBody(h, w, r, contentType)
		if err != nil {
			return
		}
	}

	if notNil(params) {
		params, err = readRequestParams(h, w, r)
		if err != nil {
			return
		}
	}

	resp, err := processRequest(w, r, handler, body, params)
	if err != nil {
		return
	}

	if notNil(resp) {
		writeResponse(h, w, r, resp)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func readRequestBody[Body, Params, Response any](h *Handler[Body, Params, Response], w http.ResponseWriter, r *http.Request, contentType string) (Body, error) {
	if contentType == "" {
		contentType = "application/json"
	}

	bodyReader, ok := h.bodyReaders[contentType]
	if !ok {
		http.Error(w, fmt.Sprintf("unsupported content type %s", contentType), http.StatusBadRequest)
		return *new(Body), fmt.Errorf("unsupported content type")
	}

	var body Body
	err := bodyReader.Reader(r, &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return *new(Body), err
	}

	return body, nil
}

func readRequestParams[Body, Params, Response any](h *Handler[Body, Params, Response], w http.ResponseWriter, r *http.Request) (Params, error) {
	var params Params
	err := h.paramsReader.Reader(r, &params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return *new(Params), err
	}

	return params, nil
}

func processRequest[Body, Params, Response any](w http.ResponseWriter, r *http.Request, handler ResponseFunc[Body, Params, Response], body Body, params Params) (Response, error) {
	resp, err := handler(r.Context(), body, params)
	if err != nil {
		switch e := err.(type) {
		case Error:
			http.Error(w, err.Error(), e.Status())
		case error:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return *new(Response), err
	}

	return resp, nil
}

func writeResponse[Body, Params, Response any](h *Handler[Body, Params, Response], w http.ResponseWriter, r *http.Request, resp Response) {
	accept := r.Header.Get("Accept")
	if accept == "" || accept == "*/*" {
		accept = "application/json"
	}

	writer, ok := h.writers[accept]
	if !ok {
		http.Error(w, fmt.Sprintf("unsupported accept type %s", accept), http.StatusBadRequest)
		return
	}

	err := writer.Write(w, r, resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
