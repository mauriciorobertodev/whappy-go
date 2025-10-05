package newsletter

import "github.com/mauriciorobertodev/whappy-go/internal/domain/events"

const (
	EventChangedPhoto = "newsletter:changed/photo"

	// newsletter:new/*
	EventNewTextMessage     events.EventName = "newsletter:new/text"     // Dispatched when a new text message is received from a newsletter
	EventNewImageMessage    events.EventName = "newsletter:new/image"    // Dispatched when a new image message is received from a newsletter
	EventNewVideoMessage    events.EventName = "newsletter:new/video"    // Dispatched when a new video message is received from a newsletter
	EventNewAudioMessage    events.EventName = "newsletter:new/audio"    // Dispatched when a new audio message is received from a newsletter
	EventNewVoiceMessage    events.EventName = "newsletter:new/voice"    // Dispatched when a new voice message is received from a newsletter
	EventNewDocumentMessage events.EventName = "newsletter:new/document" // Dispatched when a new document message is received from a newsletter
)
