package tradingaccount

import "errors"

var (
	ErrTradingAccountNotFound = errors.New("trading account not found")

	ErrInvalidLogin            = errors.New("invalid login")
	ErrInvalidUserID           = errors.New("invalid user id")
	ErrInvalidBroker           = errors.New("invalid broker")
	ErrInvalidInvestorPassword = errors.New("invalid investor password")

	ErrLoginAlreadyExists = errors.New("trading account login already exists")
	ErrUserNotFound       = errors.New("user not found")
)
