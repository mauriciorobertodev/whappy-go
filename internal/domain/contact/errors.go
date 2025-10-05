package contact

import "errors"

var (
	ErrNotFound   = errors.New("contact not found")
	ErrInvalidJID = errors.New("invalid jid")
)
