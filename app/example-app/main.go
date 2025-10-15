package main

import (
	"crypto/rand"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nick-friedrich/beesting/app/example-app/db"
	"github.com/nick-friedrich/beesting/app/example-app/handler"
	"github.com/nick-friedrich/beesting/app/example-app/pkg/config"
	"github.com/nick-friedrich/beesting/app/example-app/pkg/mail"
	"github.com/nick-friedrich/beesting/app/example-app/pkg/session"
	"github.com/nick-friedrich/beesting/app/example-app/pkg/validation"
)

func main() {

	// Initialize config
	config.InitConfig(&config.Config{
		BaseURL: "http://localhost:3000",
		EmailConfig: config.EmailConfig{
			From: "noreply@beesting.com",
			Name: "BeeSting",
		},
		AuthConfig: config.AuthConfig{
			ConfirmEmail: true,
		},
	})

	// Initialize validator singleton
	validation.InitValidator()

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

	// Initialize global session manager
	session.Default = session.NewSessionManager(queries)

	// Generate CSRF key
	csrfKey := make([]byte, 32)
	if _, err := rand.Read(csrfKey); err != nil {
		log.Fatal("Failed to generate CSRF key:", err)
	}

	// Initialize mailer
	mail.InitMailer(&mail.ConsoleAdapter{})

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes with CSRF protection
	r.Use(csrf.Protect(csrfKey, csrf.TrustedOrigins([]string{"localhost:3000"}), csrf.FieldName("_csrf")))

	// Static files (no CSRF needed)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	r.Get("/", handler.Home())
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Posts routes
	r.Route("/posts", func(r chi.Router) {
		r.Get("/", handler.ShowPosts(queries))
		r.Get("/{slug}", handler.ShowPost(queries))
		r.Get("/new", handler.CreatePostShow())
		r.Post("/new", handler.CreatePostSubmit(queries))
		r.Get("/{id}/edit", handler.EditPostShow(queries))
		r.Post("/{id}/edit", handler.EditPostSubmit(queries))
		r.Post("/{id}/delete", handler.DeletePostWeb(queries))

		r.Route("/api", func(r chi.Router) {
			r.Get("/", handler.ListPosts(queries))
			r.Post("/", handler.CreatePost(queries))
			r.Get("/{id}", handler.GetPost(queries))
			r.Put("/{id}", handler.UpdatePost(queries))
			r.Delete("/{id}", handler.DeletePost(queries))
			r.Post("/{id}/publish", handler.PublishPost(queries))
		})
	})

	// Auth routes
	r.Get("/login", handler.LoginHandler())
	r.Get("/register", handler.RegisterHandler())
	r.Post("/login", handler.LoginSubmitHandler(queries))
	r.Post("/register", handler.RegisterSubmitHandler(queries))
	r.Get("/logout", handler.LogoutHandler())
	r.Get("/verify-email", handler.VerifyEmailHandler(queries))
	r.Post("/resend-confirmation", handler.ResendConfirmationEmailHandler(queries))

	// 404 handler for unmatched routes
	r.NotFound(handler.NotFound())

	log.Println("🐝 Server running on http://localhost:3000")
	http.ListenAndServe(":3000", r)
}
