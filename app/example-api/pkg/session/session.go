package session

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/nick-friedrich/beesting/app/example-api/db"
)

// SessionData holds the session information
type SessionData struct {
	UserID   string
	Email    string
	Name     string
	LoggedIn bool
	UserRole string
}

// SessionManager handles session operations
type SessionManager struct {
	cookieName string
	maxAge     int
	queries    *db.Queries
}

// Default is the global session manager instance
var Default *SessionManager

// NewSessionManager creates a new session manager
func NewSessionManager(queries *db.Queries) *SessionManager {
	return &SessionManager{
		cookieName: "beesting_session",
		maxAge:     7 * 24 * 60 * 60, // 7 days in seconds
		queries:    queries,
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

	// Calculate expiration time
	expiresAt := time.Now().Add(time.Duration(sm.maxAge) * time.Second)

	// Store session in database
	_, err = sm.queries.CreateSession(context.Background(), db.CreateSessionParams{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return fmt.Errorf("failed to create session in database: %w", err)
	}

	// Set cookie with only the session token
	cookie := &http.Cookie{
		Name:     sm.cookieName,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   sm.maxAge,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
	return nil
}

// GetSession retrieves session data from the cookie and database
func (sm *SessionManager) GetSession(r *http.Request) (*SessionData, error) {
	cookie, err := r.Cookie(sm.cookieName)
	if err != nil {
		return &SessionData{LoggedIn: false}, nil
	}

	sessionToken := cookie.Value

	// Get session from database
	session, err := sm.queries.GetSession(context.Background(), sessionToken)
	if err != nil {
		if err == sql.ErrNoRows {
			// Session not found or expired
			return &SessionData{LoggedIn: false}, nil
		}
		return &SessionData{LoggedIn: false}, fmt.Errorf("failed to get session: %w", err)
	}

	// Update last accessed time
	_ = sm.queries.UpdateSessionAccess(context.Background(), sessionToken)

	// Get user details
	user, err := sm.queries.GetUser(context.Background(), session.UserID)
	if err != nil {
		return &SessionData{LoggedIn: false}, fmt.Errorf("failed to get user: %w", err)
	}

	return &SessionData{
		UserID:   user.ID,
		UserRole: user.Role,
		Email:    user.Email,
		Name:     user.Name,
		LoggedIn: true,
	}, nil
}

// ClearSession removes the session from database and clears the cookie
func (sm *SessionManager) ClearSession(w http.ResponseWriter, r *http.Request) error {
	// Get session token from cookie
	cookie, err := r.Cookie(sm.cookieName)
	if err == nil {
		// Delete session from database
		_ = sm.queries.DeleteSession(context.Background(), cookie.Value)
	}

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     sm.cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

// CleanupExpiredSessions removes all expired sessions from the database
func (sm *SessionManager) CleanupExpiredSessions() error {
	return sm.queries.DeleteExpiredSessions(context.Background())
}

// DeleteUserSessions removes all sessions for a specific user
func (sm *SessionManager) DeleteUserSessions(userID string) error {
	return sm.queries.DeleteUserSessions(context.Background(), userID)
}

// GetUserSessions retrieves all active sessions for a user
func (sm *SessionManager) GetUserSessions(userID string) ([]db.Session, error) {
	return sm.queries.GetUserSessions(context.Background(), userID)
}
