-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS trades_competition_login_idx
ON trades (competition_id, trading_account_login);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS trades_competition_login_idx;
-- +goose StatementEnd
