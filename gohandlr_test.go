package gohandlr

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/epentland/gohandlr/options"
)

type mockBodyReader struct {
	reader func(r *http.Request, data any) error
}

func (m *mockBodyReader) ContentType() string {
	return "application/json"
}

func (m *mockBodyReader) Reader(r *http.Request, data any) error {
	return m.reader(r, data)
}

type mockParamsReader struct {
	reader func(r *http.Request, data any) error
}

func (m *mockParamsReader) Reader(r *http.Request, data any) error {
	return m.reader(r, data)
}

type mockWriter struct {
	write func(w http.ResponseWriter, r *http.Request, data any) error
}

func (m *mockWriter) Accept() string {
	return "application/json"
}

func (m *mockWriter) Write(w http.ResponseWriter, r *http.Request, data any) error {
	return m.write(w, r, data)
}

func TestSetOptions(t *testing.T) {
	// Create a Handler instance
	handler := &Handler[string, int, bool]{
		writers:     make(map[string]options.Writer),
		bodyReaders: make(map[string]options.BodyReader),
	}

	// Define test options
	writerOpt := &mockWriter{}
	bodyReaderOpt := &mockBodyReader{}
	paramsReaderOpt := &mockParamsReader{}

	// Call SetOptions with test options
	SetOptions(handler, writerOpt, bodyReaderOpt, paramsReaderOpt)

	// Assert that the Handler instance is updated correctly
	if len(handler.writers) != 1 {
		t.Errorf("Expected 1 writer, got %d", len(handler.writers))
	}
	if _, ok := handler.writers["application/json"]; !ok {
		t.Errorf("Expected writer with accept type 'application/json'")
	}

	if len(handler.bodyReaders) != 1 {
		t.Errorf("Expected 1 body reader, got %d", len(handler.bodyReaders))
	}
	if _, ok := handler.bodyReaders["application/json"]; !ok {
		t.Errorf("Expected body reader with content type 'application/json'")
	}

	if !reflect.DeepEqual(handler.paramsReader, paramsReaderOpt) {
		t.Errorf("Expected params reader to be set")
	}
}

func TestNotNil(t *testing.T) {
	if !notNil(1) {
		t.Errorf("expected false, got true")
	}

	if notNil(Nil{}) {
		t.Errorf("expected true, got false")
	}
}

func TestNewHandler(t *testing.T) {
	// Test cases
	testCases := []struct {
		name            string
		options         []options.Options
		expectedErr     string
		expectedHandler func(t *testing.T, handler *Handler[string, int, bool])
	}{
		{
			name:        "Valid options",
			options:     []options.Options{&mockBodyReader{}, &mockParamsReader{}, &mockWriter{}},
			expectedErr: "",
			expectedHandler: func(t *testing.T, handler *Handler[string, int, bool]) {
				if len(handler.writers) != 1 {
					t.Errorf("Expected 1 writer, got: %d", len(handler.writers))
				}
				if _, ok := handler.writers["application/json"]; !ok {
					t.Errorf("Expected writer with accept type 'application/json'")
				}
				if len(handler.bodyReaders) != 1 {
					t.Errorf("Expected 1 body reader, got: %d", len(handler.bodyReaders))
				}
				if _, ok := handler.bodyReaders["application/json"]; !ok {
					t.Errorf("Expected body reader with content type 'application/json'")
				}
				if handler.paramsReader == nil {
					t.Errorf("Expected params reader to be set")
				}
			},
		},
		{
			name:            "Invalid options",
			options:         []options.Options{&mockBodyReader{}, &mockWriter{}},
			expectedErr:     "no params reader provided, please provide one",
			expectedHandler: nil,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := NewHandler[string, int, bool](tc.options...)

			if tc.expectedErr == "" {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
				if tc.expectedHandler != nil {
					tc.expectedHandler(t, handler)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error: %s, got nil", tc.expectedErr)
				} else if err.Error() != tc.expectedErr {
					t.Errorf("Expected error: %s, got: %v", tc.expectedErr, err)
				}
				if handler != nil {
					t.Errorf("Expected nil handler, got: %+v", handler)
				}
			}
		})
	}
}

func TestValidateHandler(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		body        interface{}
		params      interface{}
		response    interface{}
		options     []options.Options
		expectedErr string
	}{
		{
			name:        "Valid options",
			body:        "test",
			params:      123,
			response:    true,
			options:     []options.Options{&mockBodyReader{}, &mockParamsReader{}, &mockWriter{}},
			expectedErr: "",
		},
		{
			name:        "Missing body reader",
			body:        "test",
			params:      123,
			response:    true,
			options:     []options.Options{&mockParamsReader{}, &mockWriter{}},
			expectedErr: "no body reader provided, please provide one",
		},
		{
			name:        "Missing params reader",
			body:        "test",
			params:      123,
			response:    true,
			options:     []options.Options{&mockBodyReader{}, &mockWriter{}},
			expectedErr: "no params reader provided, please provide one",
		},
		{
			name:        "Missing writer",
			body:        "test",
			params:      123,
			response:    true,
			options:     []options.Options{&mockBodyReader{}, &mockParamsReader{}},
			expectedErr: "no writers provided, please provide one",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := &Handler[string, int, bool]{
				writers:     make(map[string]options.Writer),
				bodyReaders: make(map[string]options.BodyReader),
			}

			// Set options on the handler
			for _, opt := range tc.options {
				switch v := opt.(type) {
				case options.Writer:
					handler.writers[v.Accept()] = v
				case options.BodyReader:
					handler.bodyReaders[v.ContentType()] = v
				case options.ParamsReader:
					handler.paramsReader = v
				}
			}

			err := ValidateHandler(handler)

			if tc.expectedErr == "" {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error: %s, got nil", tc.expectedErr)
				} else if err.Error() != tc.expectedErr {
					t.Errorf("Expected error: %s, got: %v", tc.expectedErr, err)
				}
			}
		})
	}
}

func TestHandle(t *testing.T) {
	// Create a mock HttpHandler
	var httpHandlerCalled bool
	mockHttpHandler := func(path string, handler func(http.ResponseWriter, *http.Request)) {
		httpHandlerCalled = true
		// Assert the path
		if path != "/test" {
			t.Errorf("Expected path '/test', got '%s'", path)
		}
		// Create a mock request and response
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		resp := httptest.NewRecorder()
		// Call the handler
		handler(resp, req)
	}

	// Create a mock ResponseFunc
	var handlerCalled bool
	mockHandler := func(ctx context.Context, body string, params int) (bool, error) {
		handlerCalled = true
		// Assert the body and params
		if body != "test body" {
			t.Errorf("Expected body 'test body', got '%s'", body)
		}
		if params != 123 {
			t.Errorf("Expected params 123, got %d", params)
		}
		return true, nil
	}

	// Create mock options
	mockBodyReader := &mockBodyReader{
		reader: func(r *http.Request, data any) error {
			*(data.(*string)) = "test body"
			return nil
		},
	}
	mockParamsReader := &mockParamsReader{
		reader: func(r *http.Request, data any) error {
			*(data.(*int)) = 123
			return nil
		},
	}
	mockWriter := &mockWriter{
		write: func(w http.ResponseWriter, r *http.Request, data any) error {
			return nil
		},
	}
	opts := []options.Options{mockBodyReader, mockParamsReader, mockWriter}

	// Call the Handle function
	err := Handle(mockHttpHandler, "/test", mockHandler, opts...)

	// Assert the results
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !httpHandlerCalled {
		t.Error("Expected HttpHandler to be called")
	}
	if !handlerCalled {
		t.Error("Expected ResponseFunc to be called")
	}
}
