package render

import (
	"html/template"
	"io"
)

type Template struct {
	templates *template.Template
}

func NewTmpl() *Template {
	 return &Template{
	 	templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

