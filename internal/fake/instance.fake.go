package fake

import (
	"time"

	"github.com/google/uuid"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
)

type instanceFactory struct {
	prototype *instance.Instance
}

func InstanceFactory() *instanceFactory {
	return &instanceFactory{
		prototype: &instance.Instance{
			ID:     "",
			Name:   "",
			Phone:  "",
			JID:    "",
			LID:    "",
			Device: "",
			Status: instance.StatusCreated,

			LastQRCode: nil,

			LastLoginAt:     nil,
			LastConnectedAt: nil,
			BannedAt:        nil,
			BanExpiresAt:    nil,

			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
	}
}

func (f *instanceFactory) WithID(id string) *instanceFactory {
	f.prototype.ID = id
	return f
}

func (f *instanceFactory) WithName(name string) *instanceFactory {
	f.prototype.Name = name
	return f
}

func (f *instanceFactory) WithPhone(phone string) *instanceFactory {
	f.prototype.Phone = phone
	return f
}

func (f *instanceFactory) WithJID(jid string) *instanceFactory {
	f.prototype.JID = jid
	return f
}

func (f *instanceFactory) WithLID(lid string) *instanceFactory {
	f.prototype.LID = lid
	return f
}

func (f *instanceFactory) WithDevice(device string) *instanceFactory {
	f.prototype.Device = device
	return f
}

func (f *instanceFactory) WithStatus(status instance.InstanceStatus) *instanceFactory {
	f.prototype.Status = status
	return f
}

func (f *instanceFactory) Connected() *instanceFactory {
	now := time.Now().UTC()

	f.prototype.Status = instance.StatusConnected
	f.prototype.LastLoginAt = &now
	f.prototype.LastConnectedAt = &now
	f.prototype.UpdatedAt = now
	f.prototype.LastQRCode = nil
	f.prototype.BannedAt = nil
	f.prototype.BanExpiresAt = nil

	return f
}

func (f *instanceFactory) WithCreatedAt(createdAt time.Time) *instanceFactory {
	f.prototype.CreatedAt = createdAt
	return f
}

func (f *instanceFactory) Create() *instance.Instance {
	p := *f.prototype

	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}

	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = time.Now().UTC()
	}

	if p.ID == "" {
		p.ID = uuid.NewString()
	}

	return &p
}
