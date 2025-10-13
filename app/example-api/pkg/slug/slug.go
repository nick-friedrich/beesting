package slug

import (
	"regexp"
	"strings"
	"unicode"
)

// Generate creates a URL-friendly slug from a title
func Generate(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces and underscores with hyphens
	slug = regexp.MustCompile(`[\s_]+`).ReplaceAllString(slug, "-")

	// Remove all non-alphanumeric characters except hyphens
	slug = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(slug, "")

	// Remove multiple consecutive hyphens
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")

	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}

// Validate checks if a slug is valid
func Validate(slug string) error {
	if slug == "" {
		return &SlugError{Message: "slug cannot be empty"}
	}

	if len(slug) > 100 {
		return &SlugError{Message: "slug cannot be longer than 100 characters"}
	}

	// Check if slug contains only allowed characters
	if !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(slug) {
		return &SlugError{Message: "slug can only contain lowercase letters, numbers, and hyphens"}
	}

	// Check if slug starts or ends with hyphen
	if strings.HasPrefix(slug, "-") || strings.HasSuffix(slug, "-") {
		return &SlugError{Message: "slug cannot start or end with a hyphen"}
	}

	// Check if slug contains consecutive hyphens
	if strings.Contains(slug, "--") {
		return &SlugError{Message: "slug cannot contain consecutive hyphens"}
	}

	return nil
}

// SlugError represents a slug validation error
type SlugError struct {
	Message string
}

func (e *SlugError) Error() string {
	return e.Message
}

// IsASCII checks if a string contains only ASCII characters
func IsASCII(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}
