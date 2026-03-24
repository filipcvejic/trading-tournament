-- name: CreateTrackedTrade :exec
INSERT INTO tracked_trades (
    position_id,
    symbol,
    side,
    open_price,
    stop_loss,
    volume,
    opened_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
ON CONFLICT (position_id) DO NOTHING;

-- name: CloseTrackedTrade :execrows
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
    volume,
    opened_at,
    closed_at
FROM tracked_trades
ORDER BY opened_at ASC;