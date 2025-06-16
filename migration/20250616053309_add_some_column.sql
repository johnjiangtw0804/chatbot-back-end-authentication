-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN phone_number VARCHAR(20) NOT NULL DEFAULT 'N/A';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN IF EXISTS phone_number;
-- +goose StatementEnd
