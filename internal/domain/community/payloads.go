package community

import "time"

type Changer struct {
	JID string `json:"jid"`
	LID string `json:"lid"`
}

type PayloadCommunityChangedPhoto struct {
	JID       string    `json:"jid"`       // Community ID
	Changer   Changer   `json:"changer"`   // Who changed the photo
	Photo     string    `json:"photo"`     // When empty means the photo was removed
	Removed   bool      `json:"removed"`   // true if photo was removed
	Timestamp time.Time `json:"timestamp"` // When the change happened
}
