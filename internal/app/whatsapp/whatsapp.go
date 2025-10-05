package whatsapp

import (
	"context"
	"errors"
	"io"

	"github.com/mauriciorobertodev/whappy-go/internal/domain/chat"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/contact"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/file"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/group"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/message"
)

const (
	MaxPresenceDuration = 60_000 // 60 seconds
)

var (
	ErrTimeout                   = errors.New("pairing timeout")
	ErrClientOutdated            = errors.New("client outdated")
	ErrScannedWithoutMultiDevice = errors.New("scanned without multidevice")
	ErrPictureNotFound           = errors.New("picture not found")      // that user or group does not have a profile picture or invalid user/group
	ErrNoProfilePicture          = errors.New("no profile picture")     // that user or group does not have a profile picture
	ErrHiddenProfilePicture      = errors.New("hidden profile picture") // the user has hidden their profile picture from you
)

type PairingEventType int

const (
	QRCodeGenerated PairingEventType = iota
	PairingSuccess
	PairingTimeout
	PairingError
)

type PairingEvent struct {
	Type   PairingEventType
	Code   string // QR code
	Phone  string
	JID    string
	LID    string
	Device string
	Error  error
}

type Ping struct {
	IsLoggedIn  bool
	IsConnected bool
}

type MediaKind int

const (
	MediaImage MediaKind = iota
	MediaVideo
	MediaAudio
	MediaDocument
)

type PhoneStatus struct {
	Original     string  `json:"original"`
	JID          string  `json:"jid"`
	Phone        string  `json:"phone"`
	Exists       bool    `json:"exists"`
	IsBusiness   bool    `json:"is_business"`
	VerifiedName *string `json:"verified_name"`
}

type BlocklistAction string

const (
	BlocklistActionBlock   BlocklistAction = "block"
	BlocklistActionUnblock BlocklistAction = "unblock"
)

type WhatsAppGateway interface {
	// Auth
	PairingQrCode(ctx context.Context, inst *instance.Instance) (<-chan PairingEvent, error)
	Logout(ctx context.Context, instance *instance.Instance) error
	// Connection
	Connect(ctx context.Context, instance *instance.Instance) error
	Disconnect(ctx context.Context, instance *instance.Instance) error
	Ping(ctx context.Context, inst *instance.Instance) (Ping, error)

	UploadFile(ctx context.Context, inst *instance.Instance, file io.ReadCloser, kind MediaKind, mime string) (*file.File, error)

	// Messages
	GenerateMessageID(ctx context.Context, inst *instance.Instance) (string, error)
	SendTextMessage(ctx context.Context, inst *instance.Instance, message *message.Message) (*message.Message, error)
	SendImageMessage(ctx context.Context, inst *instance.Instance, message *message.Message) (*message.Message, error)
	SendVideoMessage(ctx context.Context, inst *instance.Instance, message *message.Message) (*message.Message, error)
	SendAudioMessage(ctx context.Context, inst *instance.Instance, message *message.Message) (*message.Message, error)
	SendVoiceMessage(ctx context.Context, inst *instance.Instance, message *message.Message) (*message.Message, error)
	SendDocumentMessage(ctx context.Context, inst *instance.Instance, message *message.Message) (*message.Message, error)
	SendReaction(ctx context.Context, inst *instance.Instance, message *message.Message) (*message.Message, error)
	ReadMessages(ctx context.Context, inst *instance.Instance, chat string, ids []string, sender string) error
	// Chat
	SendChatPresence(ctx context.Context, inst *instance.Instance, presence chat.Presence) error
	// Contacts
	CheckPhones(ctx context.Context, inst *instance.Instance, phones []string) ([]PhoneStatus, error)
	GetContact(ctx context.Context, inst *instance.Instance, phoneOrJID string) (*contact.Contact, error)
	GetContacts(ctx context.Context, inst *instance.Instance) ([]*contact.Contact, error)
	// Blocklist
	GetBlockList(ctx context.Context, inst *instance.Instance) ([]string, error)
	UpdateBlockList(ctx context.Context, inst *instance.Instance, phoneOrJID string, action BlocklistAction) ([]string, error)
	// Groups
	GetGroup(ctx context.Context, inst *instance.Instance, groupJID string, withParticipants *bool) (*group.Group, error)
	GetGroups(ctx context.Context, inst *instance.Instance, withParticipants *bool) ([]*group.Group, error)
	GetPictureURL(ctx context.Context, inst *instance.Instance, phoneOrJID string, preview bool, isCommunity bool) (string, error)
	JoinGroup(ctx context.Context, inst *instance.Instance, inviteCode string) (*group.Group, error)
	LeaveGroup(ctx context.Context, inst *instance.Instance, groupJID string) (*group.Group, error)
	GetGroupInviteLink(ctx context.Context, inst *instance.Instance, groupJID string, refresh bool) (string, error)
	GroupUpdatePhoto(ctx context.Context, inst *instance.Instance, groupJID string, photo *group.GroupPhoto) (string, error)
	GroupRemovePhoto(ctx context.Context, inst *instance.Instance, groupJID string) error
	GroupUpdateName(ctx context.Context, inst *instance.Instance, groupJID string, name string) error
	GroupUpdateDescription(ctx context.Context, inst *instance.Instance, groupJID string, description string) error
	GroupUpdateTopic(ctx context.Context, inst *instance.Instance, groupJID string, topic string) error
	GroupUpdateSetting(ctx context.Context, inst *instance.Instance, groupJID string, setting group.GroupSettingName, policy group.GroupSettingPolicy) error
	GroupUpdateMessageDuration(ctx context.Context, inst *instance.Instance, groupJID string, duration group.GroupMessageDuration) error
	GroupUpdateParticipants(ctx context.Context, inst *instance.Instance, groupJID string, participants []string, action group.ParticipantsAction) ([]*group.GroupParticipant, error)
	GroupCreate(ctx context.Context, inst *instance.Instance, name string, participants []string) (*group.Group, error)
}
