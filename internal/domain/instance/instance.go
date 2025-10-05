package instance

import (
	"time"

	"github.com/google/uuid"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
)

const (
	MaxNameLength = 255
)

type Instance struct {
	ID     string
	Name   string
	Phone  string
	JID    string
	LID    string
	Device string
	Status InstanceStatus

	LastQRCode *string

	LastLoginAt     *time.Time
	LastConnectedAt *time.Time
	BannedAt        *time.Time
	BanExpiresAt    *time.Time
	UpdatedAt       time.Time
	CreatedAt       time.Time
}

func NewInstance(name string) *Instance {
	id, _ := uuid.NewV7()

	return &Instance{
		ID:        id.String(),
		Name:      name,
		Phone:     "",
		JID:       "",
		Device:    "",
		Status:    StatusCreated,
		UpdatedAt: time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	}
}

// #region Rules
func (i *Instance) CanPair() error {
	if i.Status.IsBanned() {
		return ErrInstanceIsBanned
	}

	if i.Status.IsLoggedIn() {
		return ErrInstanceAlreadyLoggedIn
	}

	if i.Status.IsPairing() {
		return ErrInstanceIsPairing
	}

	return nil
}

// rule: just can log in if not banned, is pairing
func (i *Instance) CanLoginWith(phone string) error {
	if i.Status.IsBanned() {
		return ErrInstanceIsBanned
	}

	if !i.Status.IsPairing() {
		return ErrInstanceNotPairing
	}

	if i.HasLoggedInBefore() && i.Phone != phone {
		return ErrInstancePhoneConflict
	}

	return nil
}

// rule: just can log out if not banned, is logged in and not already logged out
func (i *Instance) CanLogout() error {
	if i.Status.IsLoggedOut() {
		return ErrInstanceAlreadyLoggedOut
	}

	if !i.Status.IsLoggedIn() {
		return ErrInstanceNotLoggedIn
	}

	return nil
}

// rule: just can start connection if not banned, is logged in and not already connected or connecting
func (i *Instance) CanConnect() error {
	if i.Status.IsBanned() {
		return ErrInstanceIsBanned
	}

	if i.Status.IsConnected() {
		return ErrInstanceAlreadyConnected
	}

	if i.Status.IsConnecting() {
		return ErrInstanceIsConnecting
	}

	if !i.Status.IsLoggedIn() {
		return ErrInstanceNotLoggedIn
	}

	return nil
}

// rule: just can disconnect if is connected
func (i *Instance) CanDisconnect() error {
	if !i.Status.IsConnected() {
		return ErrInstanceNotConnected
	}

	return nil
}

// rule: just can send message if not banned, is connected
func (i *Instance) CanSendMessage() error {
	if i.Status.IsBanned() {
		return ErrInstanceIsBanned
	}

	if !i.Status.IsConnected() {
		return ErrInstanceNotConnected
	}

	return nil
}

// #region Markers
func (i *Instance) MarkPairing() {
	i.Status = StatusPairing
	i.touch()
}

func (i *Instance) MarkLoggedIn(phone string, JID string, LID string, device string) {
	i.Status = StatusLoggedIn
	i.Phone = phone
	i.JID = JID
	i.LID = LID
	i.Device = device
	now := time.Now()
	i.LastLoginAt = &now
	i.touch()
}

func (i *Instance) MarkLoggedOut() {
	if i.Status.IsNew() {
		i.Status = StatusCreated
	} else {
		i.Status = StatusLoggedOut
	}

	i.touch()
}

func (i *Instance) MarkConnecting() {
	i.Status = StatusConnecting
	i.touch()
}

func (i *Instance) MarkConnected() {
	i.Status = StatusConnected
	now := time.Now()
	i.LastConnectedAt = &now
	i.touch()
}

func (i *Instance) MarkDisconnected() {
	if i.Status.IsNew() {
		i.Status = StatusCreated
	} else {
		i.Status = StatusLoggedIn
	}

	i.touch()
}

func (i *Instance) MarkPermanentlyBanned() {
	i.Status = StatusPermanentlyBanned
	now := time.Now()
	i.BannedAt = &now
	i.touch()
}

func (i *Instance) MarkTemporaryBanned(duration time.Duration) {
	i.Status = StatusTemporaryBanned
	now := time.Now()
	expiresAt := now.Add(duration)
	i.BanExpiresAt = &expiresAt
	i.BannedAt = &now
	i.touch()
}

// #region Actions
func (i *Instance) StartPairing() events.Event {
	i.MarkPairing()
	return i.EventPairingStarted()
}

func (i *Instance) Connect() events.Event {
	i.MarkConnected()
	return i.EventConnected()
}

func (i *Instance) Logout() events.Event {
	i.MarkLoggedOut()
	return i.EventLoggedOut()
}

func (i *Instance) FailPairing(code FailPairingFailedCode, attemptedPhone string, cause error) events.Event {
	i.touch()
	if i.HasLoggedInBefore() {
		i.Status = StatusLoggedOut
	} else {
		i.Status = StatusCreated
	}

	return i.EventPairingFailed(code, attemptedPhone, cause)
}

func (i *Instance) AttachQRCode(qrCode string) events.Event {
	i.LastQRCode = &qrCode
	return i.EventQRCodeGenerated(qrCode)
}

func (i *Instance) LoginWith(phone, JID, LID, device string) events.Event {
	i.MarkLoggedIn(phone, JID, LID, device)
	return i.EventLoggedIn()
}

func (i *Instance) StartConnection() events.Event {
	i.MarkConnecting()
	return i.EventConnecting()
}

func (i *Instance) Disconnect() events.Event {
	i.MarkDisconnected()
	return i.EventDisconnected()
}

// #region Helpers
func (i *Instance) touch() {
	i.UpdatedAt = time.Now().UTC()
}

func (i *Instance) HasLoggedInBefore() bool {
	return i.LastLoginAt != nil && i.Device != ""
}
