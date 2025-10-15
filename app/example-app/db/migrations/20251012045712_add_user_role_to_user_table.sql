-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN role TEXT NOT NULL DEFAULT 'user';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN role;
-- +goose StatementEnd
