package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/nick-friedrich/beesting/app/example-api/db"
	"github.com/nick-friedrich/beesting/app/example-api/pkg/config"
	"github.com/nick-friedrich/beesting/app/example-api/pkg/mail"
	passwordPkg "github.com/nick-friedrich/beesting/app/example-api/pkg/password"
	"github.com/nick-friedrich/beesting/app/example-api/pkg/session"
	"github.com/nick-friedrich/beesting/app/example-api/pkg/validation"
	"github.com/nick-friedrich/beesting/app/example-api/types"
	"github.com/nick-friedrich/beesting/app/example-api/views"
	authviews "github.com/nick-friedrich/beesting/app/example-api/views/auth"
	"github.com/oklog/ulid/v2"
)

// Legacy validation functions removed - now using validator v10 struct-based validation

func sendConfirmationEmail(user *db.User) error {
	config := config.GetConfig()
	mailer := mail.GetMailer()
	err := mailer.SendEmail(&mail.Email{
		From:    fmt.Sprintf("%s <%s>", config.EmailConfig.Name, config.EmailConfig.From),
		To:      user.Email,
		Subject: "Confirm your email",
		Body: fmt.Sprintf(`
Please click the link to confirm your email: 
%s/verify-email?token=%s
		`,
			config.BaseURL,
			user.Confirmemailtoken.String,
		),
	})
	if err != nil {
		return err
	}
	return nil
}

func LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.Default.GetSession(r)

		if sessionData.LoggedIn {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		registered := r.URL.Query().Get("registered") == "true"
		confirmEmail := r.URL.Query().Get("confirmEmail") == "true"
		emailSent := r.URL.Query().Get("emailSent") == "true"

		var successMessage string
		if registered {
			successMessage = "Registration successful! Please log in with your credentials."
		}
		if confirmEmail {
			successMessage = "Email confirmed! Please log in with your credentials."
		}
		if emailSent {
			successMessage = "Confirmation email sent! Please check your inbox and click the verification link."
		}

		views.Layout(
			authviews.Login(authviews.LoginProps{
				Errors:                types.AuthValidationErrors{},
				Email:                 "",
				SuccessMessage:        successMessage,
				ErrorMessage:          "",
				ShowResendConfirmLink: false,
				Request:               r,
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
					Errors:                errors,
					Email:                 form.Email,
					SuccessMessage:        "",
					ErrorMessage:          "",
					ShowResendConfirmLink: false,
					Request:               r,
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
					Errors:                types.AuthValidationErrors{General: "Invalid email or password"},
					Email:                 form.Email,
					SuccessMessage:        "",
					ErrorMessage:          "",
					ShowResendConfirmLink: false,
					Request:               r,
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
					Email:                 form.Email,
					SuccessMessage:        "",
					ErrorMessage:          "",
					ShowResendConfirmLink: false,
					Request:               r,
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
					Email:                 form.Email,
					SuccessMessage:        "",
					ErrorMessage:          "",
					ShowResendConfirmLink: false,
					Request:               r,
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
					Errors:                types.AuthValidationErrors{General: "Email not verified. Please check your email for a verification link."},
					Email:                 form.Email,
					SuccessMessage:        "",
					ErrorMessage:          "",
					ShowResendConfirmLink: true,
					Request:               r,
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
					Errors:                types.AuthValidationErrors{General: "Login successful but session error. Please try again."},
					Email:                 form.Email,
					SuccessMessage:        "",
					ErrorMessage:          "",
					ShowResendConfirmLink: false,
					Request:               r,
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

		// Error handling
		if err != nil {
			fmt.Printf("User creation error: %v\n", err)
			// Check if it's a unique constraint violation on email
			if strings.Contains(err.Error(), "UNIQUE constraint failed") && strings.Contains(err.Error(), "email") {
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

		// Send confirmation email
		err = sendConfirmationEmail(&user)
		if err != nil {
			fmt.Printf("Confirmation email error: %v\n", err)
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

func VerifyEmailHandler(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := q.GetByConfirmEmailToken(r.Context(), sql.NullString{String: token, Valid: true})
		if err != nil {
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}
		if user.ID == "" {
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}
		if user.Confirmedat.Valid {
			http.Error(w, "Email already verified", http.StatusBadRequest)
			return
		}
		if user.Confirmemailtokenexpiresat.Time.Before(time.Now()) {
			http.Error(w, "Token expired", http.StatusBadRequest)
			return
		}

		// Confirm the user's email
		err = q.ConfirmUserEmail(r.Context(), user.ID)
		if err != nil {
			http.Error(w, "Failed to confirm email", http.StatusInternalServerError)
			return
		}

		// Redirect to login with success message
		http.Redirect(w, r, "/login?confirmEmail=true", http.StatusSeeOther)
	}
}

func ResendConfirmationEmailHandler(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.Default.GetSession(r)

		// Don't allow resending if already logged in
		if sessionData.LoggedIn {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Get email from form
		email := strings.TrimSpace(r.Form.Get("email"))
		if email == "" {
			views.Layout(
				authviews.Login(authviews.LoginProps{
					Errors:                types.AuthValidationErrors{General: "Email is required"},
					Email:                 "",
					SuccessMessage:        "",
					ErrorMessage:          "",
					ShowResendConfirmLink: false,
					Request:               r,
				}),
				sessionData,
				"Login",
			).Render(r.Context(), w)
			return
		}

		// Get user by email
		user, err := q.GetUserByEmail(r.Context(), email)
		if err != nil {
			// Don't reveal if user exists or not - just show success message
			http.Redirect(w, r, "/login?emailSent=true", http.StatusSeeOther)
			return
		}

		// Check if already confirmed
		if user.Confirmedat.Valid {
			views.Layout(
				authviews.Login(authviews.LoginProps{
					Errors:                types.AuthValidationErrors{General: "Email already confirmed. Please log in."},
					Email:                 email,
					SuccessMessage:        "",
					ErrorMessage:          "",
					ShowResendConfirmLink: false,
					Request:               r,
				}),
				sessionData,
				"Login",
			).Render(r.Context(), w)
			return
		}

		// Generate new confirmation token
		confirmEmailToken := ulid.Make().String()
		confirmEmailTokenExpiresAt := time.Now().Add(time.Hour * 24)

		// Update user with new token
		err = q.UpdateUser(r.Context(), db.UpdateUserParams{
			ID:                         user.ID,
			Name:                       user.Name,
			Email:                      user.Email,
			PasswordHash:               user.PasswordHash,
			Confirmemailtoken:          sql.NullString{String: confirmEmailToken, Valid: true},
			Confirmemailtokenexpiresat: sql.NullTime{Time: confirmEmailTokenExpiresAt, Valid: true},
		})
		if err != nil {
			fmt.Printf("Failed to update user with new token: %v\n", err)
			views.Layout(
				authviews.Login(authviews.LoginProps{
					Errors:                types.AuthValidationErrors{General: "Failed to send confirmation email. Please try again."},
					Email:                 email,
					SuccessMessage:        "",
					ErrorMessage:          "",
					ShowResendConfirmLink: false,
					Request:               r,
				}),
				sessionData,
				"Login",
			).Render(r.Context(), w)
			return
		}

		// Update user struct with new token for email sending
		user.Confirmemailtoken = sql.NullString{String: confirmEmailToken, Valid: true}
		user.Confirmemailtokenexpiresat = sql.NullTime{Time: confirmEmailTokenExpiresAt, Valid: true}

		// Send confirmation email
		err = sendConfirmationEmail(&user)
		if err != nil {
			fmt.Printf("Confirmation email error: %v\n", err)
			// Still redirect to success to avoid revealing if user exists
		}

		// Redirect to login with success message
		http.Redirect(w, r, "/login?emailSent=true", http.StatusSeeOther)
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
