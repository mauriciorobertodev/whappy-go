package chat

import "time"

type Sender struct {
	JID string `json:"jid"`
	LID string `json:"lid"`
}

type PayloadChatChangedPresence struct {
	Chat     string           `json:"chat"`
	Sender   Sender           `json:"sender"`
	Presence ChatPresenceType `json:"presence"`
}

type PayloadChatStateRead struct {
	Chat      string    `json:"chat"`
	Timestamp time.Time `json:"timestamp"`
}

type PayloadChatStateCleared struct {
	Chat      string    `json:"chat"`
	Timestamp time.Time `json:"timestamp"`
}

type PayloadChatStateDeleted struct {
	Chat      string    `json:"chat"`
	Timestamp time.Time `json:"timestamp"`
}

type PayloadChatChangedMute struct {
	Chat      string    `json:"chat"`
	Muted     bool      `json:"muted"`
	Timestamp time.Time `json:"timestamp"`
}

type PayloadChatChangedPin struct {
	Chat      string    `json:"chat"`
	Pinned    bool      `json:"pinned"`
	Timestamp time.Time `json:"timestamp"`
}

type PayloadChatChangedArchive struct {
	Chat      string    `json:"chat"`
	Archived  bool      `json:"archived"`
	Timestamp time.Time `json:"timestamp"`
}
