package requests

import (
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http"
)

type ReadMessagesRequest struct {
	Chat   string   `json:"chat"`
	Sender string   `json:"sender"`
	IDs    []string `json:"ids"`
}

func (r *ReadMessagesRequest) Validate() *http.ErrorBag {
	var bag = http.NewErrorBag()

	if r.Chat == "" {
		bag.Add("chat", "chat is required")
	}

	if r.Sender == "" {
		bag.Add("sender", "sender is required")
	}

	if len(r.IDs) == 0 {
		bag.Add("ids", "ids are required")
	}

	return bag
}

func (r *ReadMessagesRequest) ToInput() input.ReadMessagesInput {
	return input.ReadMessagesInput{
		Chat:   r.Chat,
		Sender: r.Sender,
		IDs:    r.IDs,
	}
}

type SendReactionMessageInput struct {
	To      string `json:"to"`
	Message string `json:"message"`
	Emoji   string `json:"emoji"`
}

func (r *SendReactionMessageInput) Validate() *http.ErrorBag {
	var bag = http.NewErrorBag()

	if r.To == "" {
		bag.Add("to", "to is required")
	}

	if r.Message == "" {
		bag.Add("message", "message is required")
	}

	return bag
}

func (r *SendReactionMessageInput) ToInput() input.SendReactionInput {
	return input.SendReactionInput{
		To:      r.To,
		Message: r.Message,
		Emoji:   r.Emoji,
	}
}
