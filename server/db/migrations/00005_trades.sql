-- +goose Up
-- +goose StatementBegin
CREATE TABLE trades (
    trading_account_login BIGINT NOT NULL REFERENCES trading_accounts(login) ON DELETE CASCADE,
    competition_id UUID NOT NULL REFERENCES competitions(id) ON DELETE CASCADE,
    position_id BIGINT NOT NULL,
    symbol TEXT NOT NULL,
    side TEXT NOT NULL,
    volume NUMERIC NOT NULL,
    open_time TIMESTAMPTZ NOT NULL,
    close_time TIMESTAMPTZ NOT NULL,
    open_price NUMERIC NOT NULL,
    close_price NUMERIC NOT NULL,
    profit NUMERIC NOT NULL,
    commission NUMERIC NOT NULL,
    swap NUMERIC NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT trades_member_fkey
        FOREIGN KEY (competition_id, trading_account_login)
            REFERENCES competition_members (competition_id, trading_account_login)
            ON DELETE CASCADE,

    CONSTRAINT trades_pkey
        PRIMARY KEY (trading_account_login, position_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE trades;
-- +goose StatementEnd
