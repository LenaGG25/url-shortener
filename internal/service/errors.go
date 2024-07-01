package service

import "errors"

var (
	ErrInvalidURL  = errors.New("URL is invalid")
	ErrURLNotFound = errors.New("requested URL not found")
)
