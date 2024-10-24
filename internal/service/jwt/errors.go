package jwt

import "errors"

var (
	ErrInvalidToken     = errors.New("Invalid Token")
	ErrUnexpectedMethod = errors.New("Unexpected signing method")
)
