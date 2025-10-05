package chat

// To listen all chat events use "chat:*"
const (
	// To listen all state events use "chat:state/*"
	ChatRead    = "chat:state/read"    // Dispatched when a chat is marked as read
	ChatCleared = "chat:state/cleared" // Dispatched when a chat is cleared
	ChatDeleted = "chat:state/deleted" // Dispatched when a chat is deleted
	// To listen all changed events use "chat:changed/*"
	ChatChangedPresence = "chat:changed/presence" // Dispatched when a chat's presence is changed
	ChatChangedMute     = "chat:changed/mute"     // Dispatched when a chat is muted or unmuted
	ChatChangedPin      = "chat:changed/pin"      // Dispatched when a chat is pinned or unpinned
	ChatChangedArchive  = "chat:changed/archive"  // Dispatched when a chat is archived or unarchived
)
