package gohandlr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Config represents the configuration options
type Config struct {
	UnMarshaler     map[string]Unmarshaler
	Marshaler       map[string]Marshaler
	Validate        Validator
	ParameterReader ParameterReader
}

func (c *Config) ReadParameter(r *http.Request, v interface{}) error {
	if c.ParameterReader == nil {
		return nil
	}
	return c.ParameterReader(r, v)
}

func (c *Config) Unmarshal(r *http.Request, v interface{}) error {
	if c.UnMarshaler == nil {
		return nil
	}
	contentType := r.Header.Get("Content-Type")
	unmarshaler, ok := c.UnMarshaler[contentType]
	if !ok {
		return nil
	}
	return unmarshaler(r, v)
}

func (c *Config) Marshal(r *http.Request, w http.ResponseWriter, v interface{}) error {
	if c.Marshaler == nil {
		return nil
	}
	acceptTypes := r.Header.Get("Accept")
	acceptedTypeList := parseAcceptHeader(acceptTypes)

	for contentType, marshaler := range c.Marshaler {
		if acceptsType(acceptedTypeList, contentType) {
			w.Header().Set("Content-Type", contentType)
			return marshaler(w, v)
		}
	}
	return fmt.Errorf("can't write to any of the accepted types")
}

func parseAcceptHeader(header string) []string {
	if header == "" || header == "*/*" {
		header = "application/json"
	}
	// Split the header by commas and trim any whitespace
	parts := strings.Split(header, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
}

func acceptsType(acceptedTypes []string, mimeType string) bool {
	for _, t := range acceptedTypes {
		if t == mimeType || strings.HasPrefix(t, mimeType) {
			return true
		}
	}
	return false
}

// Option is a function that modifies the Config
type Option func(*Config)

// Reads the request
type Unmarshaler func(r *http.Request, v interface{}) error

// Writes the response
type Marshaler func(w http.ResponseWriter, v interface{}) error

// Validates the request parameters
type Validator func(v interface{}) error

// ParameterReader reads parameters from a request
type ParameterReader func(r *http.Request, v interface{}) error

// DefaultUnMarshalJSON reads the request body and unmarshals it into the .Body field of v
func DefaultUnMarshalJSON(r *http.Request, v interface{}) error {
	// Read the request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	bodyBytes = append(bodyBytes)

	// Create a new JSON decoder for the request body
	dec := json.NewDecoder(bytes.NewReader(bodyBytes))

	// Create a map to hold the wrapped body
	wrappedBody := make(map[string]json.RawMessage)

	// Unmarshal the request body into the wrapped body map
	err = dec.Decode(&wrappedBody)
	if err != nil {
		return err
	}

	// Set the "Body" field in the wrapped body map
	wrappedBody["Body"] = bodyBytes

	// Create a new JSON encoder to write the wrapped body
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	// Encode the wrapped body map
	err = enc.Encode(wrappedBody)
	if err != nil {
		return err
	}

	// Create a new JSON decoder for the wrapped body
	dec = json.NewDecoder(&buf)

	// Unmarshal the wrapped body into v
	return dec.Decode(v)
}

func DefaultMarshalJSON(w http.ResponseWriter, v interface{}) error {
	return json.NewEncoder(w).Encode(&v)
}

func WithConfig(config Config) Option {
	return func(c *Config) {
		*c = config
	}
}

// WithUnMarshaler sets the Unmarshaller in the Config
func WithUnMarshaler(contentType string, unmarshaller Unmarshaler) Option {
	return func(c *Config) {
		c.UnMarshaler[contentType] = unmarshaller
	}
}

func WithParamsReader(reader ParameterReader) Option {
	return func(c *Config) {
		c.ParameterReader = reader
	}
}

// WithMarshaler sets the Marshaller in the Config
func WithMarshaler(contentType string, marshaller Marshaler) Option {
	return func(c *Config) {
		c.Marshaler[contentType] = marshaller
	}
}

// WithValidator sets the Validator in the Config
func WithValidator(validator Validator) Option {
	return func(c *Config) {
		c.Validate = validator
	}
}

func EmptyValidator(v interface{}) error {
	return nil
}

func EmptyParameterReader(r *http.Request, v interface{}) error {
	return nil
}

var DefaultConfig = Config{
	Validate:        EmptyValidator,
	ParameterReader: EmptyParameterReader,
	UnMarshaler: map[string]Unmarshaler{
		"application/json":                  DefaultUnMarshalJSON,
		"application/x-www-form-urlencoded": DefaultUnMarshalJSON,
	},
	Marshaler: map[string]Marshaler{
		"application/json": DefaultMarshalJSON,
	},
}

func readRequest(r *http.Request, config *Config, v interface{}) error {
	// Read the request parameters
	if err := config.ReadParameter(r, v); err != nil {
		return fmt.Errorf("failed to read parameters: %w", err)
	}

	// Read the request body
	if err := config.Unmarshal(r, v); err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}

	// Validate the request data
	if err := config.Validate(v); err != nil {
		return fmt.Errorf("failed to validate request: %w", err)
	}

	return nil
}

func HandlerNoRequestNoResponse(process func(context.Context) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Process the request
		err := process(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func HandlerWithRequestNoResponse[Request any](process func(context.Context, Request) error, options ...Option) func(w http.ResponseWriter, r *http.Request) {
	config := NewConfig(options...)
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		var err error

		// Read the request
		err = readRequest(r, config, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Process the request
		err = process(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func HandlerNoRequestWithResponse[Response any](process func(context.Context) (Response, error), options ...Option) func(w http.ResponseWriter, r *http.Request) {
	config := NewConfig(options...)
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		// Process the request
		resp, err := process(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the response
		err = config.Marshal(r, w, &resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func NewConfig(options ...Option) *Config {
	config := DefaultConfig
	for _, option := range options {
		option(&config)
	}
	return &config
}

func HandlerWithRequestWithResponse[Request, Response any](process func(context.Context, Request) (Response, error), options ...Option) http.HandlerFunc {
	config := NewConfig(options...)
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		var err error

		// Read the request
		err = readRequest(r, config, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Process the request
		resp, err := process(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the response
		err = config.Marshal(r, w, &resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
