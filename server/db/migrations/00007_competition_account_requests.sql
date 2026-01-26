-- +goose Up
-- +goose StatementBegin
CREATE TABLE competition_account_requests (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    competition_id UUID NOT NULL REFERENCES competitions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (competition_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE competition_account_requests;
-- +goose StatementEnd
