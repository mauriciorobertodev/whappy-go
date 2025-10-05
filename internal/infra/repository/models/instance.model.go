package models

import (
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
)

type SQLInstance struct {
	ID     string `db:"id"`
	Name   string `db:"name"`
	Phone  string `db:"phone"`
	JID    string `db:"jid"`
	LID    string `db:"lid"`
	Device string `db:"device"`
	Status string `db:"status"`

	LastLoginAt     *time.Time `db:"last_login_at"`
	LastConnectedAt *time.Time `db:"last_connected_at"`
	BannedAt        *time.Time `db:"banned_at"`
	BanExpiresAt    *time.Time `db:"ban_expires_at"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (s *SQLInstance) ToEntity() (*instance.Instance, error) {
	var lastLoginAt, lastConnectedAt, bannedAt, banExpiresAt *time.Time
	if s.LastLoginAt != nil {
		t := s.LastLoginAt.UTC()
		lastLoginAt = &t
	}

	if s.LastConnectedAt != nil {
		t := s.LastConnectedAt.UTC()
		lastConnectedAt = &t
	}

	if s.BannedAt != nil {
		t := s.BannedAt.UTC()
		bannedAt = &t
	}

	if s.BanExpiresAt != nil {
		t := s.BanExpiresAt.UTC()
		banExpiresAt = &t
	}

	return &instance.Instance{
		ID:     s.ID,
		Name:   s.Name,
		JID:    s.JID,
		LID:    s.LID,
		Phone:  s.Phone,
		Device: s.Device,
		Status: instance.InstanceStatus(s.Status),

		LastLoginAt:     lastLoginAt,
		LastConnectedAt: lastConnectedAt,
		BannedAt:        bannedAt,
		BanExpiresAt:    banExpiresAt,

		CreatedAt: s.CreatedAt.UTC(),
		UpdatedAt: s.UpdatedAt.UTC(),
	}, nil
}

func FromInstanceEntity(inst *instance.Instance) (*SQLInstance, error) {
	var lastLoginAt, lastConnectedAt, bannedAt, banExpiresAt *time.Time
	if inst.LastLoginAt != nil {
		t := inst.LastLoginAt.UTC()
		lastLoginAt = &t
	}

	if inst.LastConnectedAt != nil {
		t := inst.LastConnectedAt.UTC()
		lastConnectedAt = &t
	}

	if inst.BannedAt != nil {
		t := inst.BannedAt.UTC()
		bannedAt = &t
	}

	if inst.BanExpiresAt != nil {
		t := inst.BanExpiresAt.UTC()
		banExpiresAt = &t
	}

	return &SQLInstance{
		ID:     inst.ID,
		Name:   inst.Name,
		Phone:  inst.Phone,
		JID:    inst.JID,
		LID:    inst.LID,
		Device: inst.Device,
		Status: string(inst.Status),

		LastLoginAt:     lastLoginAt,
		LastConnectedAt: lastConnectedAt,
		BannedAt:        bannedAt,
		BanExpiresAt:    banExpiresAt,

		CreatedAt: inst.CreatedAt.UTC(),
		UpdatedAt: inst.UpdatedAt.UTC(),
	}, nil
}
