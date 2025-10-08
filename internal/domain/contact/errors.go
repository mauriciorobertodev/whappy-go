package contact

import "errors"

var (
	ErrNotFound   = errors.New("contact not found")
	ErrInvalidJID = errors.New("invalid jid")
	EmptyPhones   = errors.New("phones list is empty")
)
