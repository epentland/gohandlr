package handle

import (
	"encoding/json"
	"html/template"
	"net/http"
)

type Writer interface {
	Write(w http.ResponseWriter, r *http.Request, buff any) error
	Accept() string
}

var _ Writer = JSONWriter{}
var _ Writer = HTMLTemplateWriter{}

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

func WithHTMLTemplateWriter(templates *template.Template, name string) HTMLTemplateWriter {
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
