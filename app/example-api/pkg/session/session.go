package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// SessionData holds the session information
type SessionData struct {
	UserID   string
	Email    string
	Name     string
	LoggedIn bool
}

// SessionManager handles session operations
type SessionManager struct {
	cookieName string
	maxAge     int
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	return &SessionManager{
		cookieName: "beesting_session",
		maxAge:     7 * 24 * 60 * 60, // 7 days in seconds
	}
}

// generateSessionID generates a secure random session ID
func generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// SetSession creates a new session and sets the cookie
func (sm *SessionManager) SetSession(w http.ResponseWriter, userID, email, name string) error {
	sessionID, err := generateSessionID()
	if err != nil {
		return fmt.Errorf("failed to generate session ID: %w", err)
	}

	// In a real application, you would store the session data in a database or Redis
	// For now, we'll encode the user data in the session ID (not secure for production)
	sessionData := fmt.Sprintf("%s|%s|%s|%s", sessionID, userID, email, name)

	cookie := &http.Cookie{
		Name:     sm.cookieName,
		Value:    sessionData,
		Path:     "/",
		MaxAge:   sm.maxAge,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
	return nil
}

// GetSession retrieves session data from the cookie
func (sm *SessionManager) GetSession(r *http.Request) (*SessionData, error) {
	cookie, err := r.Cookie(sm.cookieName)
	if err != nil {
		return &SessionData{LoggedIn: false}, nil
	}

	// Parse session data (in production, validate against stored sessions)
	// Format: sessionID|userID|email|name
	parts := splitSessionData(cookie.Value)
	if len(parts) != 4 {
		return &SessionData{LoggedIn: false}, nil
	}

	return &SessionData{
		UserID:   parts[1],
		Email:    parts[2],
		Name:     parts[3],
		LoggedIn: true,
	}, nil
}

// ClearSession removes the session cookie
func (sm *SessionManager) ClearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     sm.cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}

// splitSessionData safely splits session data
func splitSessionData(data string) []string {
	// Simple split - in production, use proper parsing/validation
	var parts []string
	var current strings.Builder
	for _, char := range data {
		if char == '|' {
			parts = append(parts, current.String())
			current.Reset()
		} else {
			current.WriteRune(char)
		}
	}
	parts = append(parts, current.String())
	return parts
}
