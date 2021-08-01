package paste

import "errors"

var (
	ErrEmptyExpireType               = errors.New("empty expire_type or expiration")
	ErrZeroExpiration                = errors.New("zero expiration")
	ErrExpirationGreaterThanMonth    = errors.New("expiration greater than a month")
	ErrExpirationGreaterThanMaxCount = errors.New("expiration greater than max count")
	ErrInvalidExpireType             = errors.New("invalid expire_type")
	ErrEmptyContent                  = errors.New("empty content")
	ErrEmptyLang                     = errors.New("empty lang")
)
