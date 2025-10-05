package group

import (
	"time"
)

const (
	MaxNameLength        = 100
	MaxDescriptionLength = 2000
	MaxTopicLength       = 2000
	MinPhotoWidth        = 192
	MaxPhotoWidth        = 640
	MinPhotoHeight       = 192
	MaxPhotoHeight       = 640
)

type AddressingMode string

const (
	AddressingModePN  AddressingMode = "pn"
	AddressingModeLID AddressingMode = "lid"
)

type GroupMemberAddMode string

const (
	GroupMemberAddModeAdmin  GroupMemberAddMode = "admin"
	GroupMemberAddModeAnyone GroupMemberAddMode = "anyone"
)

type GroupParticipant struct {
	JID           string             `json:"jid"`
	LID           string             `json:"lid"`
	Phone         string             `json:"phone"`
	Name          string             `json:"name"`
	IsMe          bool               `json:"is_me"`
	Participation GroupParticipation `json:"participation"`
}

type GroupMessageExpiration struct {
	Enabled    bool                 `json:"enabled"`
	Expiration uint32               `json:"expiration"` // Expiration time in seconds
	Duration   GroupMessageDuration `json:"duration"`
}

type GroupType string

const (
	GroupTypeCommunity    GroupType = "community"
	GroupTypeAnnouncement GroupType = "announcement"
	GroupTypeSubGroup     GroupType = "subgroup"
	GroupTypeRegular      GroupType = "regular"
)

type GroupParticipation string

const (
	GroupParticipationNone       GroupParticipation = "none"
	GroupParticipationMember     GroupParticipation = "member"
	GroupParticipationAdmin      GroupParticipation = "admin"
	GroupParticipationSuperAdmin GroupParticipation = "super_admin"
	GroupParticipationOwner      GroupParticipation = "owner"
)

type Owner struct {
	JID   string `json:"jid"`
	LID   string `json:"lid"`
	Phone string `json:"phone"`
}

type Group struct {
	JID   string    `json:"id"`    // The JID (Jabber ID) of the group
	Owner Owner     `json:"owner"` // The owner of the group
	Type  GroupType `json:"type"`  // The type of the group

	Participation GroupParticipation `json:"participation"`

	Name        string `json:"name"`        // The name of the group
	Description string `json:"description"` // The topic of the group

	// Fast boolean flags for common group settings
	Announce   bool `json:"announce"`   // Just admins can send messages
	Hidden     bool `json:"hidden"`     // This group is incognito, used for communities
	Locked     bool `json:"locked"`     // Just admins can edit group info
	Restricted bool `json:"restricted"` // Only admins can add participants
	Approval   bool `json:"approval"`   // Joining the group requires approval

	MessageExpiration GroupMessageExpiration `json:"message_expiration"` // Whether the group is ephemeral

	ParentJID string `json:"parent_jid"` // The JID of the parent group, if any

	AddressingMode AddressingMode `json:"addressing_mode"`

	CreatedAt          time.Time `json:"created_at"`
	CreatorCountryCode string    `json:"creator_country_code"`

	Participants []*GroupParticipant `json:"participants,omitempty"`
}
