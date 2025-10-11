package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/nick-friedrich/beesting/app/example-api/db"
	passwordPkg "github.com/nick-friedrich/beesting/app/example-api/pkg/password"
	"github.com/nick-friedrich/beesting/app/example-api/pkg/session"
	"github.com/nick-friedrich/beesting/app/example-api/pkg/web"
	"github.com/oklog/ulid/v2"
)

// ValidationErrors holds form validation errors
type ValidationErrors struct {
	Email    string
	Password string
	Name     string
	General  string
}

// validateEmail validates email format
func validateEmail(email string) string {
	if email == "" {
		return "Email is required"
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return "Please enter a valid email address"
	}

	return ""
}

// validatePassword validates password strength
func validatePassword(password string) string {
	if password == "" {
		return "Password is required"
	}

	if len(password) < 8 {
		return "Password must be at least 8 characters long"
	}

	return ""
}

// validateName validates name field
func validateName(name string) string {
	if name == "" {
		return "Name is required"
	}

	if len(strings.TrimSpace(name)) < 2 {
		return "Name must be at least 2 characters long"
	}

	return ""
}

func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionManager := session.NewSessionManager()
		sessionData, _ := sessionManager.GetSession(r)

		if sessionData.LoggedIn {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		registered := r.URL.Query().Get("registered") == "true"
		web.RenderWithLayoutAndSession(w, "layout.html", "templates/auth/login.html", map[string]any{
			"registered": registered,
		}, sessionData)
	}
}

func Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionManager := session.NewSessionManager()
		sessionData, _ := sessionManager.GetSession(r)

		if sessionData.LoggedIn {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		web.RenderWithLayoutAndSession(w, "layout.html", "templates/auth/register.html", map[string]any{}, sessionData)
	}
}

func LoginSubmit(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		sessionManager := session.NewSessionManager()
		sessionData, _ := sessionManager.GetSession(r)

		if sessionData.LoggedIn {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var errors ValidationErrors
		var hasErrors bool = false

		r.ParseForm()
		email := strings.TrimSpace(r.Form.Get("email"))
		password := r.Form.Get("password")

		// Validate email
		if emailErr := validateEmail(email); emailErr != "" {
			errors.Email = emailErr
			hasErrors = true
		}

		// Validate password
		if passwordErr := validatePassword(password); passwordErr != "" {
			errors.Password = passwordErr
			hasErrors = true
		}

		// If there are validation errors, show the form again
		if hasErrors {

			web.RenderWithLayout(w, "layout.html", "templates/auth/login.html", map[string]any{
				"errors": errors,
				"email":  email, // Preserve email value
			})
			return
		}

		// Authenticate user
		user, err := q.GetUserByEmail(r.Context(), email)
		if err != nil {
			// User not found or database error
			fmt.Printf("Login attempt failed: email=%s, error=%v\n", email, err)
			web.RenderWithLayout(w, "layout.html", "templates/auth/login.html", map[string]any{
				"error": "Invalid email or password",
				"email": email, // Preserve email value
			})
			return
		}

		// Verify password using Argon2
		passwordMatch, err := passwordPkg.VerifyPassword(password, user.PasswordHash)
		if err != nil {
			fmt.Printf("Password verification error: %v\n", err)
			web.RenderWithLayout(w, "layout.html", "templates/auth/login.html", map[string]any{
				"error": "Authentication error. Please try again.",
				"email": email,
			})
			return
		}

		if !passwordMatch {
			fmt.Printf("Login attempt failed: email=%s, invalid password\n", email)
			web.RenderWithLayout(w, "layout.html", "templates/auth/login.html", map[string]any{
				"error": "Invalid email or password",
				"email": email, // Preserve email value
			})
			return
		}

		// Login successful - create session
		err = sessionManager.SetSession(w, user.ID, user.Email, user.Name)
		if err != nil {
			fmt.Printf("Session creation error: %v\n", err)
			sessionManager := session.NewSessionManager()
			sessionData, _ := sessionManager.GetSession(r)

			web.RenderWithLayoutAndSession(w, "layout.html", "templates/auth/login.html", map[string]any{
				"error": "Login successful but session error. Please try again.",
				"email": email,
			}, sessionData)
			return
		}

		fmt.Printf("Login successful: user_id=%s, email=%s\n", user.ID, email)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func RegisterSubmit(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionManager := session.NewSessionManager()
		sessionData, _ := sessionManager.GetSession(r)

		if sessionData.LoggedIn {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var errors ValidationErrors
		var hasErrors bool = false

		r.ParseForm()
		name := strings.TrimSpace(r.Form.Get("name"))
		email := strings.TrimSpace(r.Form.Get("email"))
		password := r.Form.Get("password")
		confirmPassword := r.Form.Get("confirm_password")

		// Validate name
		if nameErr := validateName(name); nameErr != "" {
			errors.Name = nameErr
			hasErrors = true
		}

		// Validate email
		if emailErr := validateEmail(email); emailErr != "" {
			errors.Email = emailErr
			hasErrors = true
		}

		// Validate password
		if passwordErr := validatePassword(password); passwordErr != "" {
			errors.Password = passwordErr
			hasErrors = true
		}

		// Validate password confirmation
		if password != confirmPassword {
			errors.Password = "Passwords do not match"
			hasErrors = true
		}

		// If there are validation errors, show the form again
		if hasErrors {
			web.RenderWithLayout(w, "layout.html", "templates/auth/register.html", map[string]any{
				"errors": errors,
				"name":   name, // Preserve form values
				"email":  email,
			})
			return
		}

		// Hash the password using Argon2
		passwordHash, err := passwordPkg.HashPassword(password)
		if err != nil {
			fmt.Printf("Password hashing error: %v\n", err)
			web.RenderWithLayout(w, "layout.html", "templates/auth/register.html", map[string]any{
				"errors": ValidationErrors{
					General: "Registration error. Please try again.",
				},
				"name":  name,
				"email": email,
			})
			return
		}

		// Create user with hashed password
		user, err := q.CreateUser(r.Context(), db.CreateUserParams{
			ID:           ulid.Make().String(),
			Name:         name,
			Email:        email,
			PasswordHash: passwordHash,
		})
		if err != nil {
			fmt.Printf("User creation error: %v\n", err)
			// Check if it's a unique constraint violation
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				if strings.Contains(err.Error(), "email") {
					web.RenderWithLayout(w, "layout.html", "templates/auth/register.html", map[string]any{
						"errors": ValidationErrors{
							Email: "Email already exists. Please use a different email.",
						},
						"name":  name,
						"email": email,
					})
				} else if strings.Contains(err.Error(), "name") {
					web.RenderWithLayout(w, "layout.html", "templates/auth/register.html", map[string]any{
						"errors": ValidationErrors{
							Name: "Username already exists. Please choose a different name.",
						},
						"name":  name,
						"email": email,
					})
				} else {
					web.RenderWithLayout(w, "layout.html", "templates/auth/register.html", map[string]any{
						"errors": ValidationErrors{
							General: "Registration error. Please try again.",
						},
						"name":  name,
						"email": email,
					})
				}
			} else {
				web.RenderWithLayout(w, "layout.html", "templates/auth/register.html", map[string]any{
					"errors": ValidationErrors{
						General: "Registration error. Please try again.",
					},
					"name":  name,
					"email": email,
				})
			}
			return
		}

		fmt.Printf("User created: %+v\n", user)

		// Simulate successful registration
		http.Redirect(w, r, "/login?registered=true", http.StatusSeeOther)
	}
}

func Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionManager := session.NewSessionManager()
		sessionManager.ClearSession(w)

		fmt.Printf("User logged out\n")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
