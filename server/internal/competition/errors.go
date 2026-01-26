package competition

import "errors"

var (
	ErrCompetitionNotFound                = errors.New("competition not found")
	ErrCompetitionAlreadyStarted          = errors.New("competition already started")
	ErrAlreadyJoined                      = errors.New("trading account already joined competition")
	ErrMemberNotFound                     = errors.New("competition member not found")
	ErrInvalidAccountSize                 = errors.New("invalid account size")
	ErrInvalidCompetitionName             = errors.New("invalid competition name")
	ErrInvalidCompetitionTimeRange        = errors.New("invalid competition time range")
	ErrInvalidTradingAccountLogin         = errors.New("invalid trading account login")
	ErrInvalidPositionID                  = errors.New("invalid position id")
	ErrInvalidSymbol                      = errors.New("invalid symbol")
	ErrInvalidSide                        = errors.New("invalid side")
	ErrInvalidTradeTimeRange              = errors.New("invalid trade time range")
	ErrNotCompetitionMember               = errors.New("trading account is not a member of this competition")
	ErrTradingAccountNotFound             = errors.New("trading account not found")
	ErrAccountSizeNotSet                  = errors.New("account size is not set for this competition member")
	ErrInvalidBroker                      = errors.New("invalid broker")
	ErrInvalidInvestorPassword            = errors.New("invalid investor password")
	ErrTradingAccountAlreadyExistsForUser = errors.New("user already has trading account")
	ErrTradingAccountLoginTaken           = errors.New("trading account login already taken")
)
