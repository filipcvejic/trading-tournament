-- name: GetCompetitionLeaderboard :many
SELECT
    cm.trading_account_login::BIGINT as trading_account_login,
    
    ROW_NUMBER() OVER (
    ORDER BY
      COALESCE(
        (COALESCE(SUM(t.profit + t.commission + t.swap), 0) / NULLIF(cm.account_size, 0)) * 100,
        0
      ) DESC
  )::INT AS rank,

    u.username,
    cm.account_size::FLOAT8 AS account_size,

    COALESCE(SUM(t.profit + t.commission + t.swap), 0)::FLOAT8 AS profit,

    (cm.account_size + COALESCE(SUM(t.profit + t.commission + t.swap), 0))::FLOAT8 AS equity,

    COALESCE(
            (COALESCE(SUM(t.profit + t.commission + t.swap), 0) / NULLIF(cm.account_size, 0)) * 100,
            0
    )::FLOAT8 AS gain_percent

FROM competition_members cm
JOIN trading_accounts ta ON ta.login = cm.trading_account_login
JOIN users u ON u.id = ta.user_id
LEFT JOIN trades t ON t.trading_account_login = cm.trading_account_login
AND t.competition_id = cm.competition_id

WHERE cm.competition_id = $1

GROUP BY
    cm.trading_account_login,
    cm.account_size,
    u.username

ORDER BY gain_percent DESC
LIMIT $2 OFFSET $3;