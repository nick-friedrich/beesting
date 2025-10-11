-- +goose Up
-- +goose StatementBegin
ALTER TABLE posts ADD COLUMN slug TEXT NOT NULL DEFAULT '';
CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF NOT EXISTS idx_posts_slug;
ALTER TABLE posts DROP COLUMN slug;
-- +goose StatementEnd