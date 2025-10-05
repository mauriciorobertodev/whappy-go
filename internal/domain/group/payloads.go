package group

import "time"

type PayloadGroupPhotoChanged struct {
	JID     string `json:"jid"`
	Photo   string `json:"photo"`
	Removed bool   `json:"removed"`
}

type PayloadGroupNameChanged struct {
	JID       string    `json:"jid"`
	Name      string    `json:"name"`
	Changer   Changer   `json:"changer"`
	Timestamp time.Time `json:"timestamp"`
}

type PayloadGroupDescriptionChanged struct {
	JID         string    `json:"jid"`
	Description string    `json:"description"`
	Changer     Changer   `json:"changer"`
	Deleted     bool      `json:"deleted"`
	Timestamp   time.Time `json:"timestamp"`
}

type PayloadGroupChangedPermission struct {
	JID       string    `json:"jid"`
	Enabled   bool      `json:"enabled"`
	Changer   Changer   `json:"changer"`
	Timestamp time.Time `json:"timestamp"`
}

type PayloadGroupExpirationChanged struct {
	JID        string               `json:"id"`
	Enabled    bool                 `json:"enabled"`
	Expiration uint32               `json:"expiration"` // Expiration time in seconds
	Duration   GroupMessageDuration `json:"duration"`   // One of "off", "24h", "7d", "90d", "custom"
	Changer    Changer              `json:"changer"`
	Timestamp  time.Time            `json:"timestamp"`
}

// Novo participante entrou
// group:participants/joined
type PayloadGroupParticipantsJoined struct {
	JID          string    `json:"jid"`
	Participants []string  `json:"participants"`
	Timestamp    time.Time `json:"timestamp"`
}

// Participante saiu
// group:participants/left
type PayloadGroupParticipantsLeft struct {
	JID          string    `json:"jid"`          // identifier
	Participants []string  `json:"participants"` // identifiers
	Timestamp    time.Time `json:"timestamp"`
}

// Participante promovido a admin
// group:participants/promoted
type PayloadGroupParticipantsPromoted struct {
	JID          string    `json:"jid"`
	Participants []string  `json:"participants"`
	Promoter     Changer   `json:"promoter"`
	Timestamp    time.Time `json:"timestamp"`
}

// Admin rebaixado
// group:participants/demoted
type PayloadGroupParticipantsDemoted struct {
	JID          string    `json:"jid"`          // identifier
	Participants []string  `json:"participants"` // identifiers
	Demoter      Changer   `json:"demoter"`      // who demoted
	Timestamp    time.Time `json:"timestamp"`
}

type Changer struct {
	JID string `json:"jid"`
	LID string `json:"lid"`
}

type PayloadGroupChangedPhoto struct {
	JID       string    `json:"jid"`     // Group ID
	Changer   Changer   `json:"changer"` // Who changed the photo
	Photo     string    `json:"photo"`   // When empty means the photo was removed
	Removed   bool      `json:"removed"`
	Timestamp time.Time `json:"timestamp"`
}
