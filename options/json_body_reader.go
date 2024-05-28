package options

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WithJSONBodyReader() JSONBodyReader {
	return JSONBodyReader{}
}

type JSONBodyReader struct{}

var _ BodyReader = JSONBodyReader{}

func (j JSONBodyReader) Reader(r *http.Request, buff any) error {
	// Read the request body
	err := json.NewDecoder(r.Body).Decode(&buff)
	if err != nil {
		return fmt.Errorf("error decoding request body: %v", err)
	}

	return nil
}

func (j JSONBodyReader) ContentType() string {
	return "application/json"
}
