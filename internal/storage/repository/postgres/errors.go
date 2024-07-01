package postgres

import "errors"

var (
	ErrObjectNotFound = errors.New("object not found")
	ErrUpdateFailed   = errors.New("error to update")
)
