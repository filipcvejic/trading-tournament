package competition

import "errors"

var (
	ErrNotFound                = errors.New("not found")
	ErrAlreadyStarted          = errors.New("already started")
	ErrAlreadyJoined           = errors.New("already joined")
	ErrMemberNotFound          = errors.New("member not found")
	ErrNotMember               = errors.New("not member")
	ErrInvalidName             = errors.New("invalid name")
	ErrInvalidTimeRange        = errors.New("invalid time range")
	ErrInvalidAccountSize      = errors.New("invalid account size")
	ErrAccountSizeNotSet       = errors.New("account size not set")
	ErrInvalidLogin            = errors.New("invalid login")
	ErrInvalidPositionID       = errors.New("invalid position id")
	ErrInvalidSymbol           = errors.New("invalid symbol")
	ErrInvalidSide             = errors.New("invalid side")
	ErrInvalidTradeTimeRange   = errors.New("invalid trade time range")
	ErrInvalidBroker           = errors.New("invalid broker")
	ErrInvalidInvestorPassword = errors.New("invalid investor password")
	ErrAccountAlreadyExists    = errors.New("account already exists")
	ErrLoginTaken              = errors.New("login taken")
	ErrTradingAccountNotFound  = errors.New("trading account not found")
)
