-- +goose Up
-- +goose StatementBegin
CREATE TABLE competition_members (
    competition_id UUID NOT NULL REFERENCES competitions(id) ON DELETE CASCADE,
    trading_account_login BIGINT NOT NULL REFERENCES trading_accounts(login) ON DELETE CASCADE,
    account_size NUMERIC NOT NULL,
    
    PRIMARY KEY (competition_id, trading_account_login)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE competition_members;
-- +goose StatementEnd
