-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN confirmedAt DATETIME;
ALTER TABLE users ADD COLUMN confirmEmailToken TEXT;
ALTER TABLE users ADD COLUMN confirmEmailTokenExpiresAt DATETIME;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN confirmedAt;
ALTER TABLE users DROP COLUMN confirmEmailToken;
ALTER TABLE users DROP COLUMN confirmEmailTokenExpiresAt;
-- +goose StatementEnd
