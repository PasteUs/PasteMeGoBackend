package paste

import "errors"

var (
	ErrZeroExpireMinute               = errors.New("zero expire time")
	ErrZeroExpireCount                = errors.New("zero expire count")
	ErrExpireMinuteGreaterThanMonth   = errors.New("expire minute greater than a month")
	ErrExpireCountGreaterThanMaxCount = errors.New("expire count greater than max count")
	ErrEmptyContent                   = errors.New("empty content")
	ErrEmptyLang                      = errors.New("empty lang")
	ErrInvalidLang                    = errors.New("invalid lang")
	ErrUnencryptedPassword            = errors.New("unencrypted password ")
	ErrQueryDBFailed                  = errors.New("query from db failed")
	ErrSaveFailed                     = errors.New("save failed")
	ErrUnauthorized                   = errors.New("unauthorized")
	ErrInvalidKeyLength               = errors.New("invalid key length")
	ErrInvalidKeyFormat               = errors.New("invalid key format")
)
