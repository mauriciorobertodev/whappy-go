package chat

type ChatPresenceType string

const (
	ChatPresenceTyping    ChatPresenceType = "typing"
	ChatPresenceRecording ChatPresenceType = "recording"
	ChatPresencePaused    ChatPresenceType = "paused"
)

func (t ChatPresenceType) IsValid() bool {
	switch t {
	case ChatPresenceTyping, ChatPresenceRecording, ChatPresencePaused:
		return true
	default:
		return false
	}
}

type Presence struct {
	To   string
	Type ChatPresenceType
}
