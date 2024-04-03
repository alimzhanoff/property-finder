package storage

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrExists             = errors.New("exists")
	ErrUnexpectedRowCount = errors.New("unexpected row count")
)
