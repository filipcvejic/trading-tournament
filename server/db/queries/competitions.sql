-- name: CreateCompetition :one
INSERT INTO competitions (
    id, name, starts_at, ends_at
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetCompetitionStartTime :one
SELECT starts_at
FROM competitions
WHERE id = $1;


-- name: GetCompetitionByID :one
SELECT * FROM competitions
WHERE id = $1;

-- name: ListCompetitions :many
SELECT * FROM competitions
ORDER BY starts_at DESC;

-- name: GetCompetitionUserState :one
SELECT
    EXISTS (
        SELECT 1
        FROM competition_account_requests car
        WHERE car.user_id = $1
        AND car.competition_id = $2
    ) AS has_requested_account,
    EXISTS (
        SELECT 1
        FROM competition_members cm
        JOIN trading_accounts ta ON ta.login = cm.trading_account_login
        WHERE ta.user_id = $1
        AND cm.competition_id = $2
    ) AS has_joined;

-- name: GetCurrentCompetition :one
SELECT *
FROM competitions
WHERE now() < ends_at
ORDER BY starts_at ASC
LIMIT 1;