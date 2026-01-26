package crypto

import "errors"

var ErrInvalidKeyLength = errors.New("crypto: key must be 32 bytes (AES-256)")
