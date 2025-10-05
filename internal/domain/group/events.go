package group

import "github.com/mauriciorobertodev/whappy-go/internal/domain/events"

// To listen all group events use "group:*"
const (
	// To listen all changed events use "group:changed/*"
	EventChangedPhoto       = "group:changed/photo"       // Dispatched when a group photo is changed
	EventChangedName        = "group:changed/name"        // Dispatched when a group name is changed
	EventChangedDescription = "group:changed/description" // Dispatched when a group description is changed
	EventChangedLocked      = "group:changed/locked"      // Dispatched when a group is locked or unlocked (enabled: only admins can send messages; disabled: all participants can send messages)
	EventChangedAnnounce    = "group:changed/announce"    // Dispatched when a group announce setting is changed (enabled: only admins can edit group info; disabled: all participants can edit group info)
	EventChangedRestricted  = "group:changed/restricted"  // Dispatched when a group restrict members setting is changed (enabled: only admins can add participants; disabled: all participants can add participants)
	EventChangedApproval    = "group:changed/approval"    // Dispatched when a group approval setting is changed (enabled: new participants need admin approval to join; disabled: new participants can join without admin approval)
	EventChangedExpiration  = "group:changed/expiration"  // Dispatched when a group ephemeral setting is changed (enabled: messages will disappear after expiration time; disabled: messages will not disappear)

	// To listen all participants events use "group:participants/*"
	EventParticipantsPromoted = "group:participants/promoted" // Dispatched when an user is promoted to admin
	EventParticipantsDemoted  = "group:participants/demoted"  // Dispatched when an user is demoted to participant
	EventParticipantsJoined   = "group:participants/joined"   // Dispatched when an user is added to a group by an admin
	EventParticipantsLeft     = "group:participants/left"     // Dispatched when an user leaves a group

	// To listen all new message events, use the prefix "group:new/*"
	EventNewTextMessage     events.EventName = "group:new/text"     // Dispatched when a new text message is received from a group
	EventNewImageMessage    events.EventName = "group:new/image"    // Dispatched when a new image message is received from a group
	EventNewVideoMessage    events.EventName = "group:new/video"    // Dispatched when a new video message is received from a group
	EventNewAudioMessage    events.EventName = "group:new/audio"    // Dispatched when a new audio message is received from a group
	EventNewVoiceMessage    events.EventName = "group:new/voice"    // Dispatched when a new voice message is received from a group
	EventNewDocumentMessage events.EventName = "group:new/document" // Dispatched when a new document message is received from a group
)
