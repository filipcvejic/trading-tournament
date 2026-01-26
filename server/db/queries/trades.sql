-- name: InsertTrade :exec
INSERT INTO trades (
    trading_account_login, competition_id, position_id, symbol, side, volume, open_time, close_time, open_price, close_price, profit, commission, swap
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
) ON CONFLICT (trading_account_login, position_id)
DO NOTHING;

-- name: ListTradesByAccountLogin :many
SELECT * FROM trades
WHERE trading_account_login = $1
ORDER BY close_time DESC;