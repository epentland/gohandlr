package handle

import (
	"encoding/json"
	"html/template"
	"net/http"
)

type Writer interface {
	Write(w http.ResponseWriter, r *http.Request, data any) error
	Accept() string
}

var _ Writer = JSONWriter{}
var _ Writer = HTMLTemplateWriter{}

func NewJsonWriter() JSONWriter {
	return JSONWriter{}
}

type JSONWriter struct{}

func (j JSONWriter) Write(w http.ResponseWriter, r *http.Request, data interface{}) error {
	// Assume data is a map that can be converted to JSON
	jsonData, err := json.Marshal(data)
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

func NewHTMLTemplateWriter(templates *template.Template, name string) HTMLTemplateWriter {
	return HTMLTemplateWriter{
		Templates: templates,
		Name:      name,
	}
}

type HTMLTemplateWriter struct {
	Templates *template.Template
	Name      string
}

func (h HTMLTemplateWriter) Accept() string {
	return "text/html"
}

func (h HTMLTemplateWriter) Write(w http.ResponseWriter, r *http.Request, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	return h.Templates.ExecuteTemplate(w, h.Name, data)
}
