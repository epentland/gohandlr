package main

import (
	"context"
	"net/http"
	"strings"
)

func RegisterPostFiles(i PostFilesManager) (string, string, func(http.ResponseWriter, *http.Request)) {
	return "POST", "/files", func(w http.ResponseWriter, r *http.Request) {
		var err error

		var data User
		// Read the request body
		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case "application/json":
			data, err = i.ReadApplicationJson(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		// Process the data
		resp, err := i.Handle(r.Context(), data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the response
		acceptHeader := r.Header.Get("Accept")
		acceptedTypes := parseAcceptHeader(acceptHeader)
		switch {
		case acceptsType(acceptedTypes, "application/json"):
			err = i.Write200ApplicationJson(w, resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}

	}
}

type PostFilesManager interface {
	Handle(context.Context, User) (User, error)

	ReadApplicationJson(*http.Request) (User, error)
	//ValidateApplicationJson(*http.Request) error

	Write200ApplicationJson(http.ResponseWriter, User) error
}

type User struct {
	Name  string `json:"Name"`
	Age   int32  `json:"Age"`
	Email string `json:"Email"`
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
