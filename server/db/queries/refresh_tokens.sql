-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at)
VALUES ($1, $2, $3)
    RETURNING token, user_id, expires_at, revoked, created_at;

-- name: GetRefreshToken :one
SELECT token, user_id, expires_at, revoked, created_at
FROM refresh_tokens
WHERE token = $1
    LIMIT 1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked = TRUE
WHERE token = $1;
