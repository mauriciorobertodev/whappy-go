package events

import (
	"encoding/json"
	"strings"
	"time"
)

// Name é o identificador do tipo de evento.
type EventName string

const (
	// Eventos de mensagem
	EventMessageQueued EventName = "message.queued"
	EventMessageSent   EventName = "message.sent"
)

type Event struct {
	Name       EventName `json:"name"`
	Payload    any       `json:"payload"`
	OccurredAt time.Time `json:"occurred_at"`
	InstanceID *string   `json:"instance_id,omitempty"`
}

func New(name EventName, payload any, instanceID *string) Event {
	return Event{
		Name:       name,
		Payload:    payload,
		OccurredAt: time.Now(),
		InstanceID: instanceID,
	}
}

func (e Event) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Event) Matches(events []string) bool {
	for _, pattern := range events {
		if pattern == "*" {
			return true
		}

		if strings.HasSuffix(pattern, "/*") {
			prefix := strings.TrimSuffix(pattern, "/*")
			if strings.HasPrefix(string(e.Name), prefix+"/") {
				return true
			}
		}

		if strings.HasSuffix(pattern, ":*") {
			prefix := strings.TrimSuffix(pattern, ":*")
			if strings.HasPrefix(string(e.Name), prefix+":") {
				return true
			}
		}

		if pattern == string(e.Name) {
			return true
		}
	}
	return false
}

// message:*
// message:new/*
// message:new/text
// message:new/image
// message:new/video
// message:new/audio
// message:new/voice
// message:new/location
// message:new/contact
// message:new/sticker
// message:new/document

// message:read
// message:played
// message:delivered

// message:deleted/*
// message:deleted/text
// message:deleted/image
// message:deleted/video
// message:deleted/audio
// message:deleted/voice
// message:deleted/location
// message:deleted/contact
// message:deleted/sticker
// message:deleted/document

// group:*
// group:participants/joined
// group:participants/leaved
// group:participants/promoted
// group:participants/demoted
// group:update/photo
// group:update/name
// group:update/description

// group:update/announce
// group:update/locked
// group:update/restricted
// group:update/approval

// group:update/expiration

// session:*
// session:connected
// session:disconnected
// session:logged_in
// session:logged_out

// pair:*
// pair:success
// pair:failure
// pair:qr
