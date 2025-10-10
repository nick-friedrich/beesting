package beesting

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

type App struct {
	Tmpl *template.Template
	Mux  *http.ServeMux
}

type HandlerFunc func(http.ResponseWriter, *http.Request)

func NewApp(embedFS ...embed.FS) *App {
	var tmpl *template.Template
	if len(embedFS) > 0 {
		tmpl = template.Must(template.ParseFS(embedFS[0], "templates/*.html"))
	}
	return &App{Tmpl: tmpl, Mux: http.NewServeMux()}
}

// Handle registers a handler for the given pattern
func (a *App) Handle(pattern string, handler HandlerFunc) {
	a.Mux.HandleFunc(pattern, handler)
}

// Get registers a GET handler
func (a *App) Get(pattern string, handler HandlerFunc) {
	a.Mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	})
}

// Post registers a POST handler
func (a *App) Post(pattern string, handler HandlerFunc) {
	a.Mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	})
}

// Put registers a PUT handler
func (a *App) Put(pattern string, handler HandlerFunc) {
	a.Mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	})
}

// Delete registers a DELETE handler
func (a *App) Delete(pattern string, handler HandlerFunc) {
	a.Mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	})
}

// Patch registers a PATCH handler
func (a *App) Patch(pattern string, handler HandlerFunc) {
	a.Mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	})
}

// Static serves static files from an embed.FS or directory
func (a *App) Static(pattern string, embedFS embed.FS, dir string) {
	fileServer := http.FileServer(http.FS(embedFS))

	// Strip the pattern prefix and serve from the specified directory
	prefix := strings.TrimSuffix(pattern, "/")
	a.Mux.Handle(pattern, http.StripPrefix(prefix, fileServer))
}

// StaticDir serves static files from a filesystem directory
func (a *App) StaticDir(pattern string, dir string) {
	fileServer := http.FileServer(http.Dir(dir))
	prefix := strings.TrimSuffix(pattern, "/")
	a.Mux.Handle(pattern, http.StripPrefix(prefix, fileServer))
}

// StaticFS serves static files from an fs.FS
func (a *App) StaticFS(pattern string, fsys fs.FS) {
	fileServer := http.FileServer(http.FS(fsys))
	prefix := strings.TrimSuffix(pattern, "/")
	a.Mux.Handle(pattern, http.StripPrefix(prefix, fileServer))
}

// Run starts the HTTP server
func (a *App) Run(addr string) {
	log.Printf("🐝 BeeSting running at http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, a.Mux))
}
