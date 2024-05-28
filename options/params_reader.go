package options

import (
	"net/http"
	"reflect"
	"strconv"
)

var _ ParamsReader = DefaultParamsReader{}

type DefaultParamsReader struct{}

func WithParamsReader() DefaultParamsReader {
	return DefaultParamsReader{}
}

func (d DefaultParamsReader) Reader(r *http.Request, buff any) error {
	// Get the type of the params struct
	paramsType := reflect.TypeOf(buff)
	if paramsType.Kind() == reflect.Ptr {
		paramsType = paramsType.Elem()
	}

	// Get the value of the params struct
	paramsValue := reflect.ValueOf(buff)
	if paramsValue.Kind() == reflect.Ptr {
		paramsValue = paramsValue.Elem()
	}

	// Iterate over the fields of the params struct
	for i := 0; i < paramsType.NumField(); i++ {
		field := paramsType.Field(i)

		// Get the path tag value
		pathTag := field.Tag.Get("path")
		if pathTag != "" {
			// Get the path parameter value
			pathValue := r.PathValue(pathTag)

			// Convert the path value to the field type
			fieldValue := paramsValue.Field(i)
			switch field.Type.Kind() {
			case reflect.Int:
				intValue, _ := strconv.Atoi(pathValue)
				fieldValue.SetInt(int64(intValue))
			case reflect.String:
				fieldValue.SetString(pathValue)
			}
		}

		// Get the query tag value
		queryTag := field.Tag.Get("query")
		if queryTag != "" {
			// Get the query parameter value
			queryValue := r.URL.Query().Get(queryTag)

			// Convert the query value to the field type
			fieldValue := paramsValue.Field(i)
			switch field.Type.Kind() {
			case reflect.Int:
				intValue, _ := strconv.Atoi(queryValue)
				fieldValue.SetInt(int64(intValue))
			case reflect.String:
				fieldValue.SetString(queryValue)
			}
		}

		// Get the ctx tag value
		ctxTag := field.Tag.Get("ctx")
		if ctxTag != "" {
			// Get the context value
			ctxValue := r.Context().Value(ctxTag)

			// Convert the context value to the field type
			fieldValue := paramsValue.Field(i)
			switch field.Type.Kind() {
			case reflect.Int:
				if intValue, ok := ctxValue.(int); ok {
					fieldValue.SetInt(int64(intValue))
				}
			case reflect.String:
				if strValue, ok := ctxValue.(string); ok {
					fieldValue.SetString(strValue)
				}
			}
		}
	}

	return nil
}
