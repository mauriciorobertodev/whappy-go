package requests

import (
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http"
	"github.com/mauriciorobertodev/whappy-go/internal/utils"
)

type CreateWebhook struct {
	URL    string   `json:"url"`
	Active bool     `json:"active"`
	Events []string `json:"events"`
}

func (r *CreateWebhook) Validate() *http.ErrorBag {
	var bag = http.NewErrorBag()

	if r.Active && !utils.IsValidURL(r.URL) {
		bag.Add("url", "url is invalid")
	}

	return bag
}

func (r *CreateWebhook) ToInput() input.CreateWebhook {
	return input.CreateWebhook{
		URL:    r.URL,
		Active: r.Active,
		Events: r.Events,
	}
}

type UpdateWebhook struct {
	Active bool     `json:"active"`
	URL    string   `json:"url"`
	Events []string `json:"events"`
}

func (r *UpdateWebhook) Validate() *http.ErrorBag {
	var bag = http.NewErrorBag()

	if r.URL != "" && !utils.IsValidURL(r.URL) {
		bag.Add("url", "url is invalid")
	}

	return bag
}

func (r *UpdateWebhook) ToInput(id string) input.UpdateWebhook {
	return input.UpdateWebhook{
		ID:     id,
		Active: r.Active,
		URL:    r.URL,
		Events: r.Events,
	}
}
