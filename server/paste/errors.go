package paste

import "errors"

var (
    ErrEmptyExpireTypeOrExpiration   = errors.New("empty expire_type or expiration")
    ErrZeroExpiration                = errors.New("zero expiration")
    ErrExpirationGreaterThanMonth    = errors.New("expiration greater than a month")
    ErrExpirationGreaterThanMaxCount = errors.New("expiration greater than max count")
    ErrInvalidExpireType             = errors.New("invalid expire_type")
)
