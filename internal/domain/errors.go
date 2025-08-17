package domain

import "errors"

var (
	ErrInvalidDomain     = errors.New("invalid domain")
	ErrInvalidUUID       = errors.New("invalid UUID")
	ErrInvalidPagination = errors.New("invalid pagination parameters")
)
