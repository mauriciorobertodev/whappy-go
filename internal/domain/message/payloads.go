package message

import "time"

type PayloadMessageDelivered struct {
	Messages  []string  `json:"messages"`
	Chat      string    `json:"chat"`
	Sender    string    `json:"sender"`
	Timestamp time.Time `json:"timestamp"`
}

type PayloadMessageRead struct {
	Messages  []string  `json:"messages"`
	Chat      string    `json:"chat"`
	Sender    string    `json:"sender"`
	Timestamp time.Time `json:"timestamp"`
}

type PayloadMessagePlayed struct {
	Messages  []string  `json:"messages"`
	Chat      string    `json:"chat"`
	Sender    string    `json:"sender"`
	Timestamp time.Time `json:"timestamp"`
}

type PayloadNewMessage struct {
	Chat      string    `json:"chat"`
	Sender    Sender    `json:"sender"`
	Message   Message   `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type PayloadMessageReaction struct {
	Chat      string    `json:"chat"`    // JID do chat
	Message   string    `json:"message"` // ID da mensagem
	Sender    string    `json:"sender"`  // JID de quem reagiu
	Emoji     string    `json:"emoji"`   // Emoji usado
	Removed   bool      `json:"removed"` // Se a reação foi removida
	Timestamp time.Time `json:"timestamp"`
}
