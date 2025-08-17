package storage

import "errors"

var (
	ErrUnsupportedStorageType = errors.New("unsupported storage type")
	ErrDomainNotFound         = errors.New("domain not found")
)
