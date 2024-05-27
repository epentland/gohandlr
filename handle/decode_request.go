package handle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type RequestDecoder[T any] func(*http.Request) (T, error)

func DecodeRequestBody[T any](r *http.Request, data *T) error {
	// Check if there is supposed to be a body
	if reflect.TypeOf(data) == nil {
		return nil
	}

	// Read the request body
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		return fmt.Errorf("error decoding request body: %v", err)
	}

	return nil
}

func DecodeRequestParams[T any](r *http.Request) (T, error) {
	var params T

	// Check if there is supposed to be params
	if reflect.TypeOf(params) == nil {
		return params, nil
	}

	// Get the type of the params struct
	paramsType := reflect.TypeOf(params)

	// Iterate over the fields of the params struct
	for i := 0; i < paramsType.NumField(); i++ {
		field := paramsType.Field(i)

		// Get the path tag value
		pathTag := field.Tag.Get("path")
		if pathTag != "" {
			// Get the path parameter value
			pathValue := chi.URLParam(r, pathTag)

			// Convert the path value to the field type
			fieldValue := reflect.ValueOf(&params).Elem().Field(i)
			switch field.Type.Kind() {
			case reflect.Int:
				intValue, _ := strconv.Atoi(pathValue)
				fieldValue.SetInt(int64(intValue))
			case reflect.String:
				fieldValue.SetString(pathValue)
				// Add more cases for other supported types
			}
		}

		// Get the query tag value
		queryTag := field.Tag.Get("query")
		if queryTag != "" {
			// Get the query parameter value
			queryValue := r.URL.Query().Get(queryTag)

			// Convert the query value to the field type
			fieldValue := reflect.ValueOf(&params).Elem().Field(i)
			switch field.Type.Kind() {
			case reflect.Int:
				intValue, _ := strconv.Atoi(queryValue)
				fieldValue.SetInt(int64(intValue))
			case reflect.String:
				fieldValue.SetString(queryValue)
				// Add more cases for other supported types
			}
		}
	}

	return params, nil
}
