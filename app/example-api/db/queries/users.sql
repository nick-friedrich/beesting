-- name: CreateUser :one
INSERT INTO users (id, name, email, password_hash)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = ?
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?
LIMIT 1;
