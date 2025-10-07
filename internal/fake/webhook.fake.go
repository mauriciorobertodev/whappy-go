package fake

import (
	"time"

	"github.com/google/uuid"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/webhook"
)

type webhookFactory struct {
	prototype *webhook.Webhook
}

func WebhookFactory() *webhookFactory {
	return &webhookFactory{
		prototype: &webhook.Webhook{
			ID:         "",
			Active:     false,
			URL:        "",
			Events:     []string{},
			InstanceID: "",
			CreatedAt:  time.Time{},
			UpdatedAt:  time.Time{},
		},
	}
}

// Métodos fluentes
func (f *webhookFactory) WithID(id string) *webhookFactory {
	f.prototype.ID = id
	return f
}

func (f *webhookFactory) WithInstanceID(instanceID string) *webhookFactory {
	f.prototype.InstanceID = instanceID
	return f
}

func (f *webhookFactory) WithURL(url string) *webhookFactory {
	f.prototype.URL = url
	return f
}

func (f *webhookFactory) WithEvents(events []string) *webhookFactory {
	f.prototype.Events = events
	return f
}

func (f *webhookFactory) WithSecret(secret string) *webhookFactory {
	f.prototype.SetSecret(secret)
	return f
}

func (f *webhookFactory) Active() *webhookFactory {
	f.prototype.Active = true
	return f
}

func (f *webhookFactory) Inactive() *webhookFactory {
	f.prototype.Active = false
	return f
}

func (f *webhookFactory) WithCreatedAt(t time.Time) *webhookFactory {
	f.prototype.CreatedAt = t
	return f
}

// Criação final
func (f *webhookFactory) Create() *webhook.Webhook {
	// clona para não vazar referência
	t := *f.prototype

	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now().UTC()
	}

	if t.UpdatedAt.IsZero() {
		t.UpdatedAt = time.Now().UTC()
	}

	if t.ID == "" {
		t.ID = uuid.NewString()
	}

	if t.GetSecret() == "" {
		t.SetSecret(uuid.NewString())
	}

	if t.InstanceID == "" {
		t.InstanceID = uuid.NewString()
	}

	if t.URL == "" {
		t.URL = FakeURL()
	}

	if len(t.Events) == 0 {
		t.Events = []string{"event.a", "event.b"}
	}

	if t.Active != false {
		t.Active = true
	}

	return &t
}

func (f *webhookFactory) CreateMany(n int) []*webhook.Webhook {
	webhooks := make([]*webhook.Webhook, n)
	for i := 0; i < n; i++ {
		webhooks[i] = f.Create()
	}
	return webhooks
}
