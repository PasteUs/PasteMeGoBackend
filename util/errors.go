package util

import "errors"

var (
	ErrInvalidKeyLength = errors.New("invalid key length")
	ErrInvalidKeyFormat = errors.New("invalid key format")
)
