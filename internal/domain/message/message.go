package message

import (
	"time"
)

const (
	MaxGenerateMessageIDs = 2000
	MaxMessageTextLength  = 65_535
	MaxCaptionLength      = 1024
)

type Content interface {
	Kind() MessageKind
}

type MessageKind string

const (
	MessageKindText     MessageKind = "text"
	MessageKindImage    MessageKind = "image"
	MessageKindVideo    MessageKind = "video"
	MessageKindAudio    MessageKind = "audio"
	MessageKindVoice    MessageKind = "voice"
	MessageKindDocument MessageKind = "document"
	MessageKindReaction MessageKind = "reaction"
)

type MessageStatus string

const (
	MessageStatusFailed    MessageStatus = "failed"
	MessageStatusPending   MessageStatus = "pending"
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
	MessageStatusPlayed    MessageStatus = "played"
	MessageStatusDeleted   MessageStatus = "deleted"
)

type Sender struct {
	JID   string `json:"jid"`
	LID   string `json:"lid"`
	Phone string `json:"phone"`
	Name  string `json:"name"`
}

type Message struct {
	ID     string      `json:"id"`
	Type   MessageKind `json:"type"`
	Sender string      `json:"sender"`
	Chat   string      `json:"chat"`

	Content Content `json:"content"`

	Expiration *uint32 `json:"expiration"`

	SentAt      *time.Time `json:"sent_at"`
	DeliveredAt *time.Time `json:"delivered_at"`
	ReadAt      *time.Time `json:"read_at"`
	ExpiresAt   *time.Time `json:"expires_at"`

	IsFromMe bool `json:"is_from_me"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	InstanceID *string `json:"instance_id"`
	ExternalID *string `json:"external_id"`
}

func NewMessage(id *string, sender string, chat string, content Content, instanceID *string, expiration *uint32, me bool) *Message {
	expiresAt := (*time.Time)(nil)

	if expiration != nil {
		t := time.Now().Add(time.Duration(*expiration) * time.Second)
		expiresAt = &t
	}

	return &Message{
		Type:       content.Kind(),
		Sender:     sender,
		Chat:       chat,
		Content:    content,
		Expiration: expiration,
		ExpiresAt:  expiresAt,
		IsFromMe:   me,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		InstanceID: instanceID,
		ExternalID: id,
	}
}

func (m *Message) IsSent() bool {
	return m.SentAt != nil
}

func (m *Message) IsDelivered() bool {
	return m.DeliveredAt != nil
}

func (m *Message) IsRead() bool {
	return m.ReadAt != nil
}

func (m *Message) IsExpired() bool {
	if m.ExpiresAt == nil {
		return false
	}
	return m.ExpiresAt.Before(time.Now())
}

func (m *Message) HasMedia() bool {
	return m.Type == MessageKindImage || m.Type == MessageKindVideo || m.Type == MessageKindAudio || m.Type == MessageKindVoice || m.Type == MessageKindDocument
}

func (m *Message) IsText() bool {
	return m.Type == MessageKindText
}

func (m *Message) IsFromOther(me string) bool {
	return m.Sender != me
}

func (m *Message) AcceptContent(content Content) bool {
	return m.Type == content.Kind()
}
