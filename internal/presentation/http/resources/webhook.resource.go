package resources

import (
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/domain/webhook"
)

type WebhookResource struct {
	ID        string    `json:"id"`
	Active    bool      `json:"active"`
	URL       string    `json:"url"`
	Events    []string  `json:"events"`
	Secret    *string   `json:"secret,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func MakeWebhookResource(webhook *webhook.Webhook, secret *string) *WebhookResource {
	return &WebhookResource{
		ID:        webhook.ID,
		Active:    webhook.Active,
		URL:       webhook.URL,
		Events:    webhook.Events,
		Secret:    secret,
		UpdatedAt: webhook.UpdatedAt.UTC(),
		CreatedAt: webhook.CreatedAt.UTC(),
	}
}
