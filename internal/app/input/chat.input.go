package input

import "github.com/mauriciorobertodev/whappy-go/internal/domain/chat"

type SendChatPresenceInput struct {
	To       string                `json:"to"`
	Type     chat.ChatPresenceType `json:"type"`
	Duration *uint32               `json:"duration"` // Duration in milliseconds
	Wait     *bool                 `json:"wait"`     // Wait for the presence to be sent before returning
}

func (inp *SendChatPresenceInput) Validate() error {
	if inp.To == "" {
		return chat.ErrInvalidJID
	}

	if !inp.Type.IsValid() {
		return chat.ErrInvalidChatPresenceType
	}

	return nil
}
