package newsletter

import "time"

type Changer struct {
	JID string
	LID string
}

type PayloadNewsletterChangedPhoto struct {
	ID        string    // Newsletter ID
	Changer   Changer   // Who changed the photo
	Photo     string    // When empty means the photo was removed
	Removed   bool      // true if photo was removed
	Timestamp time.Time `json:"timestamp"`
}
