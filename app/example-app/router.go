package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/nick-friedrich/beesting/app/example-app/db"
	"github.com/nick-friedrich/beesting/app/example-app/handler"
)

func getRouter(queries *db.Queries, csrfKey []byte) *chi.Mux {
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

		// API disabled because unprotected and unversioned
		// r.Route("/api", func(r chi.Router) {
		// 	r.Get("/", handler.ListPosts(queries))
		// 	r.Post("/", handler.CreatePost(queries))
		// 	r.Get("/{id}", handler.GetPost(queries))
		// 	r.Put("/{id}", handler.UpdatePost(queries))
		// 	r.Delete("/{id}", handler.DeletePost(queries))
		// 	r.Post("/{id}/publish", handler.PublishPost(queries))
		// })
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

	return r
}
