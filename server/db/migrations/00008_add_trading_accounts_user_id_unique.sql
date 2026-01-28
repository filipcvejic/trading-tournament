-- +goose Up
-- +goose StatementBegin
ALTER TABLE trading_accounts
ADD CONSTRAINT trading_accounts_user_id_unique
UNIQUE (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE trading_accounts
DROP CONSTRAINT trading_accounts_user_id_unique;
-- +goose StatementEnd
