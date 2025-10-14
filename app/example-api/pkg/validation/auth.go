package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/nick-friedrich/beesting/app/example-api/types"
)

// LoginForm represents login form data with validation tags
type LoginForm struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

// RegisterForm represents registration form data with validation tags
type RegisterForm struct {
	Name            string `json:"name" form:"name" validate:"required,min=2"`
	Email           string `json:"email" form:"email" validate:"required,email"`
	Password        string `json:"password" form:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" validate:"required"`
}

// ValidateLoginForm validates login form data using the singleton validator
func ValidateLoginForm(form *LoginForm) error {
	if Validate == nil {
		return fmt.Errorf("validator not initialized")
	}
	return Validate.Struct(form)
}

// ValidateRegisterForm validates registration form data using the singleton validator
func ValidateRegisterForm(form *RegisterForm) error {
	if Validate == nil {
		return fmt.Errorf("validator not initialized")
	}

	if err := Validate.Struct(form); err != nil {
		return err
	}

	// Custom validation: password confirmation
	if form.Password != form.ConfirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	return nil
}

// ConvertValidationErrors converts validator errors to AuthValidationErrors
func ConvertValidationErrors(err error) types.AuthValidationErrors {
	errors := types.AuthValidationErrors{}

	if err == nil {
		return errors
	}

	// Check if it's a validator.ValidationErrors
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, validationError := range validationErrors {
			field := strings.ToLower(validationError.Field())
			message := getErrorMessage(validationError)

			switch field {
			case "name":
				errors.Name = message
			case "email":
				errors.Email = message
			case "password":
				errors.Password = message
			case "confirmpassword":
				errors.Password = message
			default:
				errors.General = message
			}
		}
	} else {
		// Handle custom errors (like password confirmation)
		errors.Password = err.Error()
	}

	return errors
}

// getErrorMessage returns a user-friendly error message for validation errors
func getErrorMessage(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return "Please enter a valid email address"
	case "min":
		if field == "Password" {
			return fmt.Sprintf("Password must be at least %s characters long", param)
		}
		return fmt.Sprintf("%s must be at least %s characters long", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, param)
	case "eqfield":
		return "Passwords do not match"
	default:
		return fmt.Sprintf("Invalid %s", strings.ToLower(field))
	}
}
