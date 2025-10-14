package validation

import (
	"github.com/go-playground/validator/v10"
)

// Global validator instance (singleton pattern)
var Validate *validator.Validate

// InitValidator initializes the global validator instance
func InitValidator() {
	Validate = validator.New()
}

// GetValidator returns the global validator instance
func GetValidator() *validator.Validate {
	return Validate
}
