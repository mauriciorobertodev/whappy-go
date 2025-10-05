package chat

import "errors"

var (
	ErrInvalidJID              = errors.New("invalid jid")
	ErrInvalidChatPresenceType = errors.New("invalid chat presence type")
)
