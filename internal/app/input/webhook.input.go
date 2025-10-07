package input

import (
	"github.com/mauriciorobertodev/whappy-go/internal/domain/webhook"
	"github.com/mauriciorobertodev/whappy-go/internal/utils"
)

type CreateWebhook struct {
	URL    string   `json:"url"`
	Active bool     `json:"active"`
	Events []string `json:"events"`
}

func (inp *CreateWebhook) Validate() error {
	if inp.Active && !utils.IsValidURL(inp.URL) {
		return webhook.ErrInvalidURL
	}

	return nil
}

type ToggleWebhook struct {
	ID     string `json:"id"`
	Active bool   `json:"active"`
}

type UpdateWebhook struct {
	ID     string   `json:"id"`
	Active bool     `json:"active"`
	URL    string   `json:"url"`
	Events []string `json:"events"`
}

func (inp *UpdateWebhook) Validate() error {
	if !utils.IsValidURL(inp.URL) {
		return webhook.ErrInvalidURL
	}

	return nil
}

type GetWebhook struct {
	ID string `json:"id"`
}

func (inp *GetWebhook) Validate() error {
	if !utils.IsUUID(inp.ID) {
		return webhook.ErrInvalidID
	}

	return nil
}

type DeleteWebhook struct {
	ID string `json:"id"`
}

func (inp *DeleteWebhook) Validate() error {
	if !utils.IsUUID(inp.ID) {
		return webhook.ErrInvalidID
	}

	return nil
}

type RenewWebhookSecret struct {
	ID string `json:"id"`
}

func (inp *RenewWebhookSecret) Validate() error {
	if !utils.IsUUID(inp.ID) {
		return webhook.ErrInvalidID
	}

	return nil
}
