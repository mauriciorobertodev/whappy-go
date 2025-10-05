package message

import "github.com/mauriciorobertodev/whappy-go/internal/domain/events"

// To listen all message events use "message:*"
const (
	// To listen all status change events, use the prefix "message:status/*"
	EventMessageSent      events.EventName = "message:status/sent"
	EventMessageRead      events.EventName = "message:status/read"
	EventMessageDelivered events.EventName = "message:status/delivered"
	EventMessagePlayed    events.EventName = "message:status/played"
	EventMessageEdited    events.EventName = "message:status/edited"
	EventMessageDeleted   events.EventName = "message:status/deleted"

	EventMessageReactionNew     events.EventName = "message:reaction/new"
	EventMessageReactionRemoved events.EventName = "message:reaction/removed"
)
