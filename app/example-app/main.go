package main

import (
	"crypto/rand"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nick-friedrich/beesting/app/example-app/db"
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

	r := getRouter(queries, csrfKey)

	log.Println("🐝 Server running on http://localhost:3000")
	http.ListenAndServe(":3000", r)
}
