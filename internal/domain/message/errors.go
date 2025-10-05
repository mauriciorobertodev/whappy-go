package message

import "errors"

var (
	ErrInvalidJID       = errors.New("invalid jid")
	ErrInvalidMessageID = errors.New("invalid message id")
	ErrInvalidQuantity  = errors.New("invalid quantity")
	ErrEmptyText        = errors.New("text cannot be empty")
	ErrInvalidURL       = errors.New("invalid url")
	ErrTextTooLong      = errors.New("text is too long")
	ErrCaptionTooLong   = errors.New("caption is too long")
	ErrImageRequired    = errors.New("image is required")
	ErrVideoRequired    = errors.New("video is required")
	ErrAudioRequired    = errors.New("audio is required")
	ErrDocumentRequired = errors.New("document is required")
)
