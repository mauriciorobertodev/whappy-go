package objects

import (
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
)

type Ping struct {
	Status          instance.InstanceStatus
	IsLoggedIn      bool
	IsConnected     bool
	LastLoginAt     *time.Time
	LastConnectedAt *time.Time
	BannedAt        *time.Time
	BanExpiresAt    *time.Time
}
