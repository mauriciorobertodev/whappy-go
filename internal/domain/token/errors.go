package token

import "errors"

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrNotFound           = errors.New("token not found")
	ErrInvalidID          = errors.New("invalid token ID")
	ErrTokenAlreadyExists = errors.New("token already exists")
)
