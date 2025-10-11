package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
)

// Template cache
var (
	templates *template.Template
	once      sync.Once
)

// LoadTemplates initializes the template cache once
func LoadTemplates() {
	once.Do(func() {
		templates = template.Must(template.New("").ParseGlob("templates/*.html"))
		templates = template.Must(templates.ParseGlob("templates/partials/*.html"))

		// Debug: list loaded templates
		log.Printf("Loaded templates: %v", templates.DefinedTemplates())
	})
}

// Render renders a template with the given data
func Render(w http.ResponseWriter, name string, data any) {
	if templates == nil {
		LoadTemplates()
	}

	log.Printf("Rendering template: %s", name)
	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
	}
}
