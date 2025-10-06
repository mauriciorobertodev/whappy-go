package models

import (
	"encoding/json"
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/domain/webhook"
)

type SQLWebhook struct {
	ID         string    `db:"id"`
	Secret     string    `db:"secret"`
	Events     string    `db:"events"` // <- agora texto
	URL        string    `db:"url"`
	Active     bool      `db:"active"`
	InstanceID string    `db:"instance_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func (s *SQLWebhook) ToEntity() *webhook.Webhook {
	var events []string
	_ = json.Unmarshal([]byte(s.Events), &events) // tolerante a erro de parse

	w := webhook.Webhook{
		ID:         s.ID,
		Events:     events,
		URL:        s.URL,
		Active:     s.Active,
		InstanceID: s.InstanceID,
		CreatedAt:  s.CreatedAt.UTC(),
		UpdatedAt:  s.UpdatedAt.UTC(),
	}
	w.SetSecret(s.Secret)
	return &w
}

func FromWebhookEntity(ent *webhook.Webhook) (*SQLWebhook, error) {
	data, err := json.Marshal(ent.Events)
	if err != nil {
		return nil, err
	}

	return &SQLWebhook{
		ID:         ent.ID,
		Secret:     ent.GetSecret(),
		Events:     string(data), // serialize
		URL:        ent.URL,
		Active:     ent.Active,
		InstanceID: ent.InstanceID,
		CreatedAt:  ent.CreatedAt.UTC(),
		UpdatedAt:  ent.UpdatedAt.UTC(),
	}, nil
}
