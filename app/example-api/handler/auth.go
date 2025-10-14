package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/nick-friedrich/beesting/app/example-api/db"
	"github.com/nick-friedrich/beesting/app/example-api/pkg/config"
	passwordPkg "github.com/nick-friedrich/beesting/app/example-api/pkg/password"
	"github.com/nick-friedrich/beesting/app/example-api/pkg/session"
	"github.com/nick-friedrich/beesting/app/example-api/pkg/validation"
	"github.com/nick-friedrich/beesting/app/example-api/types"
	"github.com/nick-friedrich/beesting/app/example-api/views"
	authviews "github.com/nick-friedrich/beesting/app/example-api/views/auth"
	"github.com/oklog/ulid/v2"
)

// Legacy validation functions removed - now using validator v10 struct-based validation

// TODO: Send confirmation email function

func LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.Default.GetSession(r)

		if sessionData.LoggedIn {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		registered := r.URL.Query().Get("registered") == "true"
		confirmEmail := r.URL.Query().Get("confirmEmail") == "true"

		var successMessage string
		if registered {
			successMessage = "Registration successful! Please log in with your credentials."
		}
		if confirmEmail {
			successMessage = "Email confirmed! Please log in with your credentials."
		}

		views.Layout(
			authviews.Login(authviews.LoginProps{
				Errors:         types.AuthValidationErrors{},
				Email:          "",
				SuccessMessage: successMessage,
				ErrorMessage:   "",
				Request:        r,
			}),
			sessionData,
			"Login",
		).Render(r.Context(), w)
	}
}

func RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.Default.GetSession(r)

		if sessionData.LoggedIn {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		views.Layout(
			authviews.Register(authviews.RegisterProps{
				Errors:  types.AuthValidationErrors{},
				Name:    "",
				Email:   "",
				Request: r,
			}),
			sessionData,
			"Register",
		).Render(r.Context(), w)

		// web.RenderWithLayoutAndSession(w, "layout.html", "templates/auth/register.html", map[string]any{}, sessionData)
	}
}

func LoginSubmitHandler(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.Default.GetSession(r)

		if sessionData.LoggedIn {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Parse form data into struct
		form := &validation.LoginForm{
			Email:    strings.TrimSpace(r.Form.Get("email")),
			Password: r.Form.Get("password"),
		}

		// Validate using struct tags
		if err := validation.ValidateLoginForm(form); err != nil {
			errors := validation.ConvertValidationErrors(err)
			views.Layout(
				authviews.Login(authviews.LoginProps{
					Errors:         errors,
					Email:          form.Email,
					SuccessMessage: "",
					ErrorMessage:   "",
					Request:        r,
				}),
				sessionData,
				"Login",
			).Render(r.Context(), w)
			return
		}

		// Authenticate user
		user, err := q.GetUserByEmail(r.Context(), form.Email)
		if err != nil {
			// User not found or database error
			fmt.Printf("Login attempt failed: email=%s, error=%v\n", form.Email, err)
			views.Layout(
				authviews.Login(authviews.LoginProps{
					Errors:         types.AuthValidationErrors{General: "Invalid email or password"},
					Email:          form.Email,
					SuccessMessage: "",
					ErrorMessage:   "",
					Request:        r,
				}),
				sessionData,
				"Login",
			).Render(r.Context(), w)

			return
		}

		// Get config and check if verified if enabled
		config := config.GetConfig()
		if config.AuthConfig.ConfirmEmail && !user.Confirmedat.Valid {
			views.Layout(
				authviews.Login(authviews.LoginProps{
					Errors:         types.AuthValidationErrors{General: "Email not verified. Please check your email for a verification link."},
					Email:          form.Email,
					SuccessMessage: "",
					ErrorMessage:   "",
					Request:        r,
				}),
				sessionData,
				"Login",
			).Render(r.Context(), w)
			return
		}

		// Verify password using Argon2
		passwordMatch, err := passwordPkg.VerifyPassword(form.Password, user.PasswordHash)
		if err != nil {
			fmt.Printf("Password verification error: %v\n", err)
			views.Layout(
				authviews.Login(authviews.LoginProps{
					Errors: types.AuthValidationErrors{
						General: "Authentication error. Please try again.",
					},
					Email:          form.Email,
					SuccessMessage: "",
					ErrorMessage:   "",
					Request:        r,
				}),
				sessionData,
				"Login",
			).Render(r.Context(), w)

			return
		}

		if !passwordMatch {
			fmt.Printf("Login attempt failed: email=%s, invalid password\n", form.Email)
			views.Layout(
				authviews.Login(authviews.LoginProps{
					Errors: types.AuthValidationErrors{
						General: "Invalid email or password",
					},
					Email:          form.Email,
					SuccessMessage: "",
					ErrorMessage:   "",
					Request:        r,
				}),
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
				authviews.Login(authviews.LoginProps{
					Errors:         types.AuthValidationErrors{General: "Login successful but session error. Please try again."},
					Email:          form.Email,
					SuccessMessage: "",
					ErrorMessage:   "",
					Request:        r,
				}),
				sessionData,
				"Login",
			).Render(r.Context(), w)

			return
		}

		fmt.Printf("Login successful: user_id=%s, email=%s\n", user.ID, form.Email)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func RegisterSubmitHandler(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.Default.GetSession(r)

		if sessionData.LoggedIn {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Parse form data into struct
		form := &validation.RegisterForm{
			Name:            strings.TrimSpace(r.Form.Get("name")),
			Email:           strings.TrimSpace(r.Form.Get("email")),
			Password:        r.Form.Get("password"),
			ConfirmPassword: r.Form.Get("confirm_password"),
		}

		// Validate using struct tags
		if err := validation.ValidateRegisterForm(form); err != nil {
			errors := validation.ConvertValidationErrors(err)
			views.Layout(
				authviews.Register(authviews.RegisterProps{
					Errors:  errors,
					Name:    form.Name,
					Email:   form.Email,
					Request: r,
				}),
				sessionData,
				"Register",
			).Render(r.Context(), w)
			return
		}

		// Hash the password using Argon2
		passwordHash, err := passwordPkg.HashPassword(form.Password)
		if err != nil {
			fmt.Printf("Password hashing error: %v\n", err)
			views.Layout(
				authviews.Register(authviews.RegisterProps{
					Errors:  types.AuthValidationErrors{General: "Registration error. Please try again."},
					Name:    form.Name,
					Email:   form.Email,
					Request: r,
				}),
				sessionData,
				"Register",
			).Render(r.Context(), w)
			return
		}

		var confirmEmailToken string
		var confirmEmailTokenExpiresAt time.Time
		config := config.GetConfig()
		if config.AuthConfig.ConfirmEmail {
			confirmEmailToken = ulid.Make().String()
			confirmEmailTokenExpiresAt = time.Now().Add(time.Hour * 24)
		}

		// Create user with hashed password
		user, err := q.CreateUser(r.Context(), db.CreateUserParams{
			ID:                         ulid.Make().String(),
			Name:                       form.Name,
			Email:                      form.Email,
			PasswordHash:               passwordHash,
			Confirmemailtoken:          sql.NullString{String: confirmEmailToken, Valid: true},
			Confirmemailtokenexpiresat: sql.NullTime{Time: confirmEmailTokenExpiresAt, Valid: true},
		})

		if err != nil {
			fmt.Printf("User creation error: %v\n", err)
			// Check if it's a unique constraint violation
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				if strings.Contains(err.Error(), "email") {
					views.Layout(
						authviews.Register(authviews.RegisterProps{
							Errors:  types.AuthValidationErrors{Email: "Email already exists. Please use a different email."},
							Name:    form.Name,
							Email:   form.Email,
							Request: r,
						}),
						sessionData,
						"Register",
					).Render(r.Context(), w)
				} else if strings.Contains(err.Error(), "name") {
					views.Layout(
						authviews.Register(authviews.RegisterProps{
							Errors:  types.AuthValidationErrors{Name: "Username already exists. Please choose a different name."},
							Name:    form.Name,
							Email:   form.Email,
							Request: r,
						}),
						sessionData,
						"Register",
					).Render(r.Context(), w)
				} else {
					views.Layout(
						authviews.Register(authviews.RegisterProps{
							Errors:  types.AuthValidationErrors{General: "Registration error. Please try again."},
							Name:    form.Name,
							Email:   form.Email,
							Request: r,
						}),
						sessionData,
						"Register",
					).Render(r.Context(), w)
				}
			} else {
				views.Layout(
					authviews.Register(authviews.RegisterProps{
						Errors:  types.AuthValidationErrors{General: "Registration error. Please try again."},
						Name:    form.Name,
						Email:   form.Email,
						Request: r,
					}),
					sessionData,
					"Register",
				).Render(r.Context(), w)
			}
			return
		}

		fmt.Printf("User created: %+v\n", user)

		// Build URL with config
		if config.AuthConfig.ConfirmEmail {
			http.Redirect(w, r, "/login?registered=true&confirmEmail=true", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/login?registered=true", http.StatusSeeOther)
		}
	}
}

func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := session.Default.ClearSession(w, r)
		if err != nil {
			fmt.Printf("Logout error: %v\n", err)
		}

		fmt.Printf("User logged out\n")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
