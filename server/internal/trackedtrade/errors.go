package trackedtrade

import "errors"

var (
	ErrInvalidEventType    = errors.New("invalid tracked trade event type")
	ErrMissingOpenFields   = errors.New("missing required fields for OPEN event")
	ErrMissingCloseFields  = errors.New("missing required fields for CLOSE event")
	ErrMissingUpdateFields = errors.New("missing required fields for UPDATE event")
	ErrTradeNotFound       = errors.New("tracked trade not found")
)
