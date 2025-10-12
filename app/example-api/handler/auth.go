package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/nick-friedrich/beesting/app/example-api/db"
	passwordPkg "github.com/nick-friedrich/beesting/app/example-api/pkg/password"
	"github.com/nick-friedrich/beesting/app/example-api/pkg/session"
	"github.com/nick-friedrich/beesting/app/example-api/types"
	"github.com/nick-friedrich/beesting/app/example-api/views"
	authviews "github.com/nick-friedrich/beesting/app/example-api/views/auth"
	"github.com/oklog/ulid/v2"
)

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
		sessionData, _ := session.Default.GetSession(r)

		if sessionData.LoggedIn {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		registered := r.URL.Query().Get("registered") == "true"

		views.Layout(
			authviews.Login(types.AuthValidationErrors{}, "", registered, r),
			sessionData,
			"Login",
		).Render(r.Context(), w)

		// web.RenderWithLayoutAndSession(w, "layout.html", "templates/auth/login.html", map[string]any{
		// 	"registered": registered,
		// }, sessionData)
	}
}

func Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.Default.GetSession(r)

		if sessionData.LoggedIn {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		views.Layout(
			authviews.Register(types.AuthValidationErrors{}, "", "", r),
			sessionData,
			"Register",
		).Render(r.Context(), w)

		// web.RenderWithLayoutAndSession(w, "layout.html", "templates/auth/register.html", map[string]any{}, sessionData)
	}
}

func LoginSubmit(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.Default.GetSession(r)

		if sessionData.LoggedIn {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var errors types.AuthValidationErrors
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
			views.Layout(
				authviews.Login(errors, email, false, r),
				sessionData,
				"Login",
			).Render(r.Context(), w)
			return
		}

		// Authenticate user
		user, err := q.GetUserByEmail(r.Context(), email)
		if err != nil {
			// User not found or database error
			fmt.Printf("Login attempt failed: email=%s, error=%v\n", email, err)
			views.Layout(
				authviews.Login(types.AuthValidationErrors{General: "Invalid email or password"}, email, false, r),
				sessionData,
				"Login",
			).Render(r.Context(), w)

			return
		}

		// Verify password using Argon2
		passwordMatch, err := passwordPkg.VerifyPassword(password, user.PasswordHash)
		if err != nil {
			fmt.Printf("Password verification error: %v\n", err)
			views.Layout(
				authviews.Login(types.AuthValidationErrors{General: "Authentication error. Please try again."}, email, false, r),
				sessionData,
				"Login",
			).Render(r.Context(), w)

			return
		}

		if !passwordMatch {
			fmt.Printf("Login attempt failed: email=%s, invalid password\n", email)
			views.Layout(
				authviews.Login(types.AuthValidationErrors{General: "Invalid email or password"}, email, false, r),
				sessionData,
				"Login",
			).Render(r.Context(), w)

			return
		}

		// Login successful - create session
		err = session.Default.SetSession(w, user.ID, user.Email, user.Name)
		if err != nil {
			fmt.Printf("Session creation error: %v\n", err)
			sessionData, _ := session.Default.GetSession(r)

			views.Layout(
				authviews.Login(types.AuthValidationErrors{General: "Login successful but session error. Please try again."}, email, false, r),
				sessionData,
				"Login",
			).Render(r.Context(), w)

			return
		}

		fmt.Printf("Login successful: user_id=%s, email=%s\n", user.ID, email)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func RegisterSubmit(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.Default.GetSession(r)

		if sessionData.LoggedIn {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var errors types.AuthValidationErrors
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
			views.Layout(
				authviews.Register(errors, name, email, r),
				sessionData,
				"Register",
			).Render(r.Context(), w)
			return
		}

		// Hash the password using Argon2
		passwordHash, err := passwordPkg.HashPassword(password)
		if err != nil {
			fmt.Printf("Password hashing error: %v\n", err)
			views.Layout(
				authviews.Register(types.AuthValidationErrors{General: "Registration error. Please try again."}, name, email, r),
				sessionData,
				"Register",
			).Render(r.Context(), w)
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
					views.Layout(
						authviews.Register(types.AuthValidationErrors{Email: "Email already exists. Please use a different email."}, name, email, r),
						sessionData,
						"Register",
					).Render(r.Context(), w)
				} else if strings.Contains(err.Error(), "name") {
					views.Layout(
						authviews.Register(types.AuthValidationErrors{Name: "Username already exists. Please choose a different name."}, name, email, r),
						sessionData,
						"Register",
					).Render(r.Context(), w)
				} else {
					views.Layout(
						authviews.Register(types.AuthValidationErrors{General: "Registration error. Please try again."}, name, email, r),
						sessionData,
						"Register",
					).Render(r.Context(), w)
				}
			} else {
				views.Layout(
					authviews.Register(types.AuthValidationErrors{General: "Registration error. Please try again."}, name, email, r),
					sessionData,
					"Register",
				).Render(r.Context(), w)
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
		err := session.Default.ClearSession(w, r)
		if err != nil {
			fmt.Printf("Logout error: %v\n", err)
		}

		fmt.Printf("User logged out\n")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
