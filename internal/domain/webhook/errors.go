package webhook

import "errors"

var (
	ErrNotFound           = errors.New("webhook not found")
	ErrInvalidURL         = errors.New("invalid webhook url")
	ErrInvalidID          = errors.New("invalid webhook id")
	ErrMaxWebhooksReached = errors.New("maximum number of webhooks reached")
)
