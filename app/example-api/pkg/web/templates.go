package web

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
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
		templates = template.New("")

		// Walk the templates directory recursively
		err := filepath.WalkDir("templates", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// Only process .html files
			if !d.IsDir() && strings.HasSuffix(path, ".html") {
				_, err := templates.ParseFiles(path)
				if err != nil {
					log.Printf("Error parsing template %s: %v", path, err)
					return err
				}
			}

			return nil
		})

		if err != nil {
			log.Fatalf("Error loading templates: %v", err)
		}

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
