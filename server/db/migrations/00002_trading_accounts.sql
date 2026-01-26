-- +goose Up
-- +goose StatementBegin
CREATE TABLE trading_accounts (
    login BIGINT PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    broker TEXT NOT NULL,
    investor_password_encrypted TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE trading_accounts;
-- +goose StatementEnd
