package user

import "github.com/mauriciorobertodev/whappy-go/internal/domain/events"

// To listen all user events use "user:*"
const (
	// To listen all change events, use the prefix "user:changed/*"
	EventChangedStatus   = "user:changed/status"   // Dispatched when the user's status message is changed
	EventChangedPhoto    = "user:changed/photo"    // Dispatched when the user's profile photo is changed
	EventChangedPresence = "user:changed/presence" // Dispatched when the user's presence (online/offline) is changed

	// To listen all new message events, use the prefix "user:new/*"
	EventNewTextMessage     events.EventName = "user:new/text"     // Dispatched when a new text message is received from a user
	EventNewImageMessage    events.EventName = "user:new/image"    // Dispatched when a new image message is received from a user
	EventNewVideoMessage    events.EventName = "user:new/video"    // Dispatched when a new video message is received from a user
	EventNewAudioMessage    events.EventName = "user:new/audio"    // Dispatched when a new audio message is received from a user
	EventNewVoiceMessage    events.EventName = "user:new/voice"    // Dispatched when a new voice message is received from a user
	EventNewDocumentMessage events.EventName = "user:new/document" // Dispatched when a new document message is received from a user
)
