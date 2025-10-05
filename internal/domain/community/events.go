package community

import "github.com/mauriciorobertodev/whappy-go/internal/domain/events"

const (
	EventChangedPhoto = "community:changed/photo"

	// community:new/*
	EventNewTextMessage     events.EventName = "community:new/text"     // Dispatched when a new text message is received from a community
	EventNewImageMessage    events.EventName = "community:new/image"    // Dispatched when a new image message is received from a community
	EventNewVideoMessage    events.EventName = "community:new/video"    // Dispatched when a new video message is received from a community
	EventNewAudioMessage    events.EventName = "community:new/audio"    // Dispatched when a new audio message is received from a community
	EventNewVoiceMessage    events.EventName = "community:new/voice"    // Dispatched when a new voice message is received from a community
	EventNewDocumentMessage events.EventName = "community:new/document" // Dispatched when a new document message is received from a community
)
