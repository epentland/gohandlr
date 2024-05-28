package options

import (
	"net/http"
	"text/template"
)

var _ Writer = HTMLTemplateWriter{}

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
