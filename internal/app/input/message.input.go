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
		return message.ErrInvalidJID
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

type SendReactionInput struct {
	To      string `json:"to"`
	Message string `json:"message"`
	Emoji   string `json:"emoji"`
}

func (inp *SendReactionInput) Validate() error {
	if inp.To == "" {
		return message.ErrInvalidJID
	}

	if inp.Message == "" {
		return message.ErrInvalidMessageID
	}

	return nil
}
