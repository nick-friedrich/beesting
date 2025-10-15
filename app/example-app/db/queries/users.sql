-- name: CreateUser :one
INSERT INTO users (id, name, email, password_hash, confirmEmailToken, confirmEmailTokenExpiresAt)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
SET name = ?,
    email = ?,
    password_hash = ?,
    confirmEmailToken = ?,
    confirmEmailTokenExpiresAt = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING id, name, email, role, confirmedat, confirmemailtoken, confirmemailtokenexpiresat, created_at, updated_at;

-- name: GetUser :one
SELECT * FROM users
WHERE id = ?
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?
LIMIT 1;

-- name: GetByConfirmEmailToken :one
SELECT * FROM users
WHERE confirmEmailToken = ?
LIMIT 1;

-- name: ConfirmUserEmail :exec
UPDATE users
SET confirmedat = CURRENT_TIMESTAMP,
    confirmEmailToken = NULL,
    confirmEmailTokenExpiresAt = NULL,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;
