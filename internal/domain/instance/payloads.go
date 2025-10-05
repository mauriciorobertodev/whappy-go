package instance

import "time"

type PayloadInstanceCreated struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type PayloadInstanceConnected struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type PayloadInstanceQRCodeGenerated struct {
	ID     string `json:"id"`
	QRCode string `json:"qr_code"`
}

type PayloadInstanceLoggedIn struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	JID    string `json:"jid"`
	Device string `json:"device"`
}

type FailPairingFailedCode string

const (
	FailPairingUnknownCode            FailPairingFailedCode = "UNKNOWN"
	FailPairingTimeoutCode            FailPairingFailedCode = "TIMEOUT"
	FailPairingConflictCode           FailPairingFailedCode = "CONFLICT"
	FailPairingClientOutdatedCode     FailPairingFailedCode = "CLIENT_OUTDATED"
	FailPairingWithoutMultideviceCode FailPairingFailedCode = "WITHOUT_MULTI_DEVICE"
)

type PayloadInstancePairingFailed struct {
	ID             string                `json:"id"`
	Code           FailPairingFailedCode `json:"code"`
	Phone          string                `json:"phone,omitempty"`
	AttemptedPhone string                `json:"attempted_phone,omitempty"`
	Error          error                 `json:"error,omitempty"`
}

type PayloadInstancePairingStarted struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PayloadInstanceConnecting struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PayloadInstanceConnectionFailed struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Error string `json:"error"`
}

type PayloadInstanceLoggedOut struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PayloadInstanceDisconnected struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PayloadInstanceBanned struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	BannedAt  time.Time  `json:"banned_at"`
	BanExpiry *time.Time `json:"ban_expires_at,omitempty"`
	Reason    *string    `json:"reason,omitempty"`
}

type PayloadInstanceGroupJoined struct {
	GroupJID          string `json:"group_jid"`
	GroupName         string `json:"group_name"`
	TotalParticipants uint16 `json:"total_participants"`
}

type PayloadInstanceGroupLeft struct {
	GroupJID  string `json:"group_jid"`
	GroupName string `json:"group_name"`
}

type PayloadInstanceBlocklistUpdated struct {
	Blocklist []string `json:"blocklist"`
}

type PayloadInstanceToken struct {
	ID     string `json:"id"`
	Token  string `json:"token"`
	Masked bool   `json:"masked"`
}
