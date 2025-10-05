package status

import "github.com/mauriciorobertodev/whappy-go/internal/domain/events"

// To listen all status events use "status:*"
const (
	// To listen all new status events use "status:new/*"
	EventStatusNewText  events.EventName = "status:new/text"  // Dispatched when a new status text is received
	EventStatusNewImage events.EventName = "status:new/image" // Dispatched when a new status image is received
	EventStatusNewVideo events.EventName = "status:new/video" // Dispatched when a new status video is received
	EventStatusNewVoice events.EventName = "status:new/voice" // Dispatched when a new status voice note is received
)
