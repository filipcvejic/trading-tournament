-- name: CreateTradingAccount :one
INSERT INTO trading_accounts (
    login, user_id, broker, investor_password_encrypted                              
) VALUES (
    $1, $2, $3, $4          
) RETURNING login, user_id, broker, created_at;

-- name: GetTradingAccountByLogin :one
SELECT login, user_id, broker, created_at
FROM trading_accounts
WHERE login = $1
LIMIT 1;
      
-- name: GetTradingAccountByUserID :one
SELECT login, user_id, broker, investor_password_encrypted, created_at
FROM trading_accounts
WHERE user_id = $1
LIMIT 1;
