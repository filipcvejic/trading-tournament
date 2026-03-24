-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user';

ALTER TABLE users
ADD CONSTRAINT users_role_check
CHECK (role IN ('user', 'admin'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP CONSTRAINT IF EXISTS users_role_check;

ALTER TABLE users
DROP COLUMN IF EXISTS role;
-- +goose StatementEnd
