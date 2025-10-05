package resources

import (
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http"
)

type InstanceResource struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	Phone           string        `json:"phone"`
	JID             string        `json:"jid"`
	LID             string        `json:"lid"`
	Device          string        `json:"device"`
	Status          string        `json:"status"`
	QRCode          *string       `json:"qr_code"`
	Reason          *string       `json:"reason"`
	LastConnectedAt *http.ISOTime `json:"last_connected_at"`
	LastLoginAt     *http.ISOTime `json:"last_login_at"`
	BannedAt        *http.ISOTime `json:"banned_at"`
	BanExpireAt     *http.ISOTime `json:"ban_expire_at"`
	UpdatedAt       http.ISOTime  `json:"updated_at"`
	CreatedAt       http.ISOTime  `json:"created_at"`
}

func MakeInstanceResource(inst *instance.Instance) *InstanceResource {
	return &InstanceResource{
		ID:              inst.ID,
		Name:            inst.Name,
		Phone:           inst.Phone,
		JID:             inst.JID,
		LID:             inst.LID,
		Device:          inst.Device,
		Status:          string(inst.Status),
		LastConnectedAt: http.NewISOTime(inst.LastConnectedAt),
		LastLoginAt:     http.NewISOTime(inst.LastLoginAt),
		BannedAt:        http.NewISOTime(inst.BannedAt),
		BanExpireAt:     http.NewISOTime(inst.BanExpiresAt),
		CreatedAt:       *http.NewISOTime(&inst.CreatedAt),
		UpdatedAt:       *http.NewISOTime(&inst.UpdatedAt),
	}
}

func MakeInstanceResources(instances []*instance.Instance) []*InstanceResource {
	resources := make([]*InstanceResource, 0, len(instances))
	for _, inst := range instances {
		resources = append(resources, MakeInstanceResource(inst))
	}

	return resources
}
