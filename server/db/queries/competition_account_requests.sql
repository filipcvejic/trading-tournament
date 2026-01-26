-- name: CreateCompetitionAccountRequest :exec
INSERT INTO competition_account_requests (
    user_id, competition_id
) VALUES ($1, $2)
ON CONFLICT (user_id, competition_id) DO NOTHING;