package input

import "github.com/mauriciorobertodev/whappy-go/internal/domain/message"

type ReadMessagesInput struct {
	Chat   string   `json:"chat"`
	Sender string   `json:"sender"`
	IDs    []string `json:"ids"`
}

func (inp *ReadMessagesInput) Validate() error {
	if inp.Chat == "" {
		return message.ErrInvalidJID
	}

	if inp.Sender == "" {
		return message.ErrInvalidJID
	}

	if len(inp.IDs) == 0 {
		return message.ErrEmptyMessageIDs
	}

	return nil
}

type GenerateMessageIDs struct {
	Quantity int `json:"quantity"`
}

func (inp *GenerateMessageIDs) Validate() error {
	if inp.Quantity <= 0 || inp.Quantity > message.MaxGenerateMessageIDs {
		return message.ErrInvalidQuantity
	}

	return nil
}
