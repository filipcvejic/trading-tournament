-- name: CreateUser :one
INSERT INTO users (
     email, username, discord_username, password_hash                   
) VALUES (
    $1, $2, $3, $4
) RETURNING id, email, username, discord_username, created_at, updated_at;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUsernameByTradingAccountLogin :one
SELECT u.username
FROM trading_accounts ta
JOIN users u ON u.id = ta.user_id
WHERE ta.login = $1
LIMIT 1;

