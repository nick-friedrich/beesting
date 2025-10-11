-- name: CreateSession :one
INSERT INTO sessions (id, user_id, expires_at)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = ?
AND expires_at > CURRENT_TIMESTAMP
LIMIT 1;

-- name: GetUserSessions :many
SELECT * FROM sessions
WHERE user_id = ?
AND expires_at > CURRENT_TIMESTAMP
ORDER BY last_accessed_at DESC;

-- name: UpdateSessionAccess :exec
UPDATE sessions
SET last_accessed_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = ?;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at <= CURRENT_TIMESTAMP;

-- name: DeleteUserSessions :exec
DELETE FROM sessions
WHERE user_id = ?;

