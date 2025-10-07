package cache

import (
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/domain/webhook"
)

type CachedWebhook struct {
	ID         string    `json:"id"`
	Active     bool      `json:"active"`
	URL        string    `json:"url"`
	Events     []string  `json:"events"`
	Secret     string    `json:"secret"`
	InstanceID string    `json:"instance_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ToCachedWebhook(w *webhook.Webhook) CachedWebhook {
	return CachedWebhook{
		ID:         w.ID,
		Active:     w.Active,
		URL:        w.URL,
		Events:     w.Events,
		Secret:     w.GetSecret(),
		InstanceID: w.InstanceID,
		CreatedAt:  w.CreatedAt,
		UpdatedAt:  w.UpdatedAt,
	}
}

func FromCachedWebhook(cw *CachedWebhook) *webhook.Webhook {
	w := &webhook.Webhook{
		ID:         cw.ID,
		Active:     cw.Active,
		URL:        cw.URL,
		Events:     cw.Events,
		InstanceID: cw.InstanceID,
		CreatedAt:  cw.CreatedAt,
		UpdatedAt:  cw.UpdatedAt,
	}
	w.SetSecret(cw.Secret)
	return w
}
