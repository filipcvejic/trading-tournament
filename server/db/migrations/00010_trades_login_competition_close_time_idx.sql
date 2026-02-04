-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS trades_login_competition_close_time_idx
ON trades (trading_account_login, competition_id, close_time);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS trades_login_competition_close_time_idx;
-- +goose StatementEnd
