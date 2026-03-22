-- name: CreateTrackedTrade :exec
INSERT INTO tracked_trades (
    position_id,
    symbol,
    side,
    open_price,
    stop_loss,
    opened_at
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: CloseTrackedTrade :exec
UPDATE tracked_trades
SET closed_at = $2
WHERE position_id = $1;

-- name: ListTrackedTrades :many
SELECT
    position_id,
    symbol,
    side,
    open_price,
    stop_loss,
    opened_at,
    closed_at
FROM tracked_trades
ORDER BY opened_at ASC;