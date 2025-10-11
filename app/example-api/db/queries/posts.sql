-- name: CreatePost :one
INSERT INTO posts (title, content, author, published)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetPost :one
SELECT * FROM posts
WHERE id = ?
LIMIT 1;

-- name: ListPosts :many
SELECT * FROM posts
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListPublishedPosts :many
SELECT * FROM posts
WHERE published = 1
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: UpdatePost :one
UPDATE posts
SET title = ?,
    content = ?,
    author = ?,
    published = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = ?;

-- name: PublishPost :exec
UPDATE posts
SET published = 1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: UnpublishPost :exec
UPDATE posts
SET published = 0,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: CountPosts :one
SELECT COUNT(*) FROM posts;

-- name: CountPublishedPosts :one
SELECT COUNT(*) FROM posts
WHERE published = 1;

