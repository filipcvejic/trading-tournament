-- +goose Up
-- +goose StatementBegin
CREATE TABLE tracked_trades (
    position_id BIGINT PRIMARY KEY,
    symbol VARCHAR(32) NOT NULL,
    side VARCHAR(10) NOT NULL CHECK (side IN ('BUY', 'SELL')),
    open_price NUMERIC(20,10) NOT NULL,
    volume NUMERIC(20,4) NOT NULL,
    stop_loss NUMERIC(20,10),
    opened_at TIMESTAMPTZ NOT NULL,
    closed_at TIMESTAMPTZ,

    CONSTRAINT tracked_trades_closed_at_valid
        CHECK (closed_at IS NULL OR closed_at >= opened_at)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tracked_trades;
-- +goose StatementEnd
