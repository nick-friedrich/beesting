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

// RenderWithLayout renders a page template within a layout
func RenderWithLayout(w http.ResponseWriter, layoutName string, pagePath string, data any) {
	RenderWithLayoutAndSession(w, layoutName, pagePath, data, nil)
}

// RenderWithLayoutAndSession renders a page template within a layout with session data
func RenderWithLayoutAndSession(w http.ResponseWriter, layoutName string, pagePath string, data any, sessionData any) {
	// Parse layout and partials first
	layoutTemplate, err := template.ParseFiles("templates/layout.html", "templates/partials/header.html", "templates/partials/footer.html")
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		return
	}

	// Clone and add the page template (which defines the content blocks)
	pageTemplate, err := template.Must(layoutTemplate.Clone()).ParseFiles(pagePath)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		return
	}

	// Merge data with session data
	templateData := make(map[string]any)
	if data != nil {
		if dataMap, ok := data.(map[string]any); ok {
			for k, v := range dataMap {
				templateData[k] = v
			}
		}
	}
	if sessionData != nil {
		templateData["Session"] = sessionData
	}

	log.Printf("Rendering template: %s with layout: %s", pagePath, layoutName)
	err = pageTemplate.ExecuteTemplate(w, layoutName, templateData)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
	}
}
