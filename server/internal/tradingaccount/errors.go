package tradingaccount

import "errors"

var (
	ErrNotFound                = errors.New("trading account not found")
	ErrInvalidLogin            = errors.New("invalid login")
	ErrInvalidBroker           = errors.New("invalid broker")
	ErrInvalidInvestorPassword = errors.New("invalid investor password")
	ErrLoginTaken              = errors.New("trading account login already taken")
)
