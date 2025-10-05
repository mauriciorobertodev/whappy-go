package user

import "time"

type PayloadUserChangedStatus struct {
	JID       string    `json:"jid"`
	LID       string    `json:"lid"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type PayloadUserChangedPhoto struct {
	JID       string    `json:"jid"`
	LID       string    `json:"lid"`
	Photo     string    `json:"photo"`   // When empty means the photo was removed
	Removed   bool      `json:"removed"` // true if photo was removed
	Timestamp time.Time `json:"timestamp"`
}

type PayloadUserPresence struct {
	JID      string     `json:"jid"`
	LID      string     `json:"lid"`
	Online   bool       `json:"online"` // true if online, false if offline
	LastSeen *time.Time `json:"last_seen"`
}
