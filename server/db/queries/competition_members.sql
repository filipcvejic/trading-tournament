-- name: JoinCompetitionBeforeStart :one
INSERT INTO competition_members (
    competition_id, trading_account_login, account_size
)
SELECT
    $1, $2, 0
FROM competitions c
WHERE c.id = $1
AND now() < c.starts_at
RETURNING competition_id;

-- name: UpdateCompetitionMemberAccountSize :exec
UPDATE competition_members
SET account_size = $3
WHERE competition_id = $1
AND trading_account_login = $2;


-- name: GetCompetitionMemberAccountSize :one
SELECT account_size
FROM competition_members
WHERE competition_id = $1
AND trading_account_login = $2;
