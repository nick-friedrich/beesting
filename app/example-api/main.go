package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nick-friedrich/beesting/app/example-api/db"
	"github.com/nick-friedrich/beesting/app/example-api/handler"
	"github.com/nick-friedrich/beesting/app/example-api/pkg/web"
)

func main() {
	// Load templates
	web.LoadTemplates()

	// Initialize database
	database, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Run migrations
	if err := db.RunMigrations(database); err != nil {
		log.Fatal(err)
	}

	// Create queries instance
	queries := db.New(database)

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Routes
	r.Get("/", handler.Home())

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Posts routes
	r.Route("/posts", func(r chi.Router) {
		r.Get("/", handler.ListPosts(queries))
		r.Post("/", handler.CreatePost(queries))
		r.Get("/{id}", handler.GetPost(queries))
		r.Put("/{id}", handler.UpdatePost(queries))
		r.Delete("/{id}", handler.DeletePost(queries))
		r.Post("/{id}/publish", handler.PublishPost(queries))
	})

	log.Println("🐝 Server running on http://localhost:3000")
	http.ListenAndServe(":3000", r)
}
