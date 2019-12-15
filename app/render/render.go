package render

import (
	"fmt"
	"html/template"
	"net/http"
)

var templates *template.Template

func Load(dir string) {
	templates = template.Must(template.ParseGlob(dir))
}

func Page(w http.ResponseWriter, page string, data interface{}) {
	templates.ExecuteTemplate(w, fmt.Sprintf("%s.html", page), data)
}
