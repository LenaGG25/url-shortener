package service

import "errors"

var (
	ErrInvalidURL  = errors.New("URL is invalid")
	ErrURLNotFound = errors.New("URL not found")
)
