package options

import (
	"encoding/json"
	"net/http"
)

var _ Writer = JSONWriter{}

func WithJsonWriter() JSONWriter {
	return JSONWriter{}
}

type JSONWriter struct{}

func (j JSONWriter) Write(w http.ResponseWriter, r *http.Request, buff interface{}) error {
	jsonData, err := json.Marshal(buff)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

	return nil
}

func (j JSONWriter) Accept() string {
	return "application/json"
}
