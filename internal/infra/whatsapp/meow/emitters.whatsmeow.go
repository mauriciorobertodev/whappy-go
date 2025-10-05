package meow

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/mauriciorobertodev/whappy-go/internal/domain/blocklist"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/chat"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/community"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/file"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/group"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/message"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/newsletter"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/privacy"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/user"
	"go.mau.fi/whatsmeow/types"
	meowEvents "go.mau.fi/whatsmeow/types/events"
)

// #region Message State Emitters
func (g *WhatsmeowGateway) emitMessageDelivered(inst *instance.Instance, evt *meowEvents.Receipt) {
	g.eventbus.Publish(events.New(
		message.EventMessageDelivered,
		message.PayloadMessageDelivered{
			Messages:  evt.MessageIDs,
			Chat:      evt.Chat.String(),
			Sender:    evt.Sender.String(),
			Timestamp: evt.Timestamp,
		},
		&inst.ID,
	))
}

func (g *WhatsmeowGateway) emitMessageRead(inst *instance.Instance, evt *meowEvents.Receipt) {
	g.eventbus.Publish(events.New(
		message.EventMessageRead,
		message.PayloadMessageRead{
			Messages:  evt.MessageIDs,
			Chat:      evt.Chat.String(),
			Sender:    evt.Sender.String(),
			Timestamp: evt.Timestamp,
		},
		&inst.ID,
	))
}

func (g *WhatsmeowGateway) emitMessagePlayed(inst *instance.Instance, evt *meowEvents.Receipt) {
	g.eventbus.Publish(events.New(
		message.EventMessagePlayed,
		message.PayloadMessagePlayed{
			Messages:  evt.MessageIDs,
			Chat:      evt.Chat.String(),
			Sender:    evt.Sender.String(),
			Timestamp: evt.Timestamp,
		},
		&inst.ID,
	))
}

// #region Status Event Emitters
func (gate *WhatsmeowGateway) emitStatusNew(inst *instance.Instance, evt *meowEvents.Message) {
	// TODO: Implement status handling if needed
	fmt.Println("Received a status")
}

// #region Privacy Emitters
func (g *WhatsmeowGateway) emitPrivacyChanged(inst *instance.Instance, evt *meowEvents.PrivacySettings) {
	if evt.CallAddChanged {
		g.eventbus.Publish(events.New(
			privacy.EventChangedCallAdd,
			privacy.PayloadPrivacyChanged{
				CallAdd: (*privacy.PrivacySetting)(&evt.NewSettings.CallAdd),
			},
			&inst.ID,
		))
	}

	if evt.GroupAddChanged {
		g.eventbus.Publish(events.New(
			privacy.EventChangedGroupAdd,
			privacy.PayloadPrivacyChanged{
				GroupAdd: (*privacy.PrivacySetting)(&evt.NewSettings.GroupAdd),
			},
			&inst.ID,
		))
	}

	if evt.LastSeenChanged {
		g.eventbus.Publish(events.New(
			privacy.EventChangedLastSeen,
			privacy.PayloadPrivacyChanged{
				LastSeen: (*privacy.PrivacySetting)(&evt.NewSettings.LastSeen),
			},
			&inst.ID,
		))
	}

	if evt.OnlineChanged {
		g.eventbus.Publish(events.New(
			privacy.EventChangedOnline,
			privacy.PayloadPrivacyChanged{
				Online: (*privacy.PrivacySetting)(&evt.NewSettings.Online),
			},
			&inst.ID,
		))
	}

	if evt.ProfileChanged {
		g.eventbus.Publish(events.New(
			privacy.EventChangedProfile,
			privacy.PayloadPrivacyChanged{
				Profile: (*privacy.PrivacySetting)(&evt.NewSettings.Profile),
			},
			&inst.ID,
		))
	}

	if evt.ReadReceiptsChanged {
		g.eventbus.Publish(events.New(
			privacy.EventChangedReadReceipts,
			privacy.PayloadPrivacyChanged{
				ReadReceipts: (*privacy.PrivacySetting)(&evt.NewSettings.ReadReceipts),
			},
			&inst.ID,
		))
	}

	if evt.StatusChanged {
		g.eventbus.Publish(events.New(
			privacy.EventChangedStatus,
			privacy.PayloadPrivacyChanged{
				Status: (*privacy.PrivacySetting)(&evt.NewSettings.Status),
			},
			&inst.ID,
		))
	}
}

// #region User Event Handlers
func (g *WhatsmeowGateway) emitUserChangedPhoto(inst *instance.Instance, evt *meowEvents.Picture) {
	jid, lid := GetJIDAndLID(evt.JID)
	g.eventbus.Publish(events.New(
		user.EventChangedPhoto,
		user.PayloadUserChangedPhoto{
			JID:       jid,
			LID:       lid,
			Photo:     evt.PictureID,
			Removed:   evt.Remove,
			Timestamp: evt.Timestamp,
		},
		&inst.ID,
	))
}

func (g *WhatsmeowGateway) emitUserChangedStatus(inst *instance.Instance, evt *meowEvents.UserAbout) {
	jid, lid := GetJIDAndLID(evt.JID)
	g.eventbus.Publish(events.New(
		user.EventChangedStatus,
		user.PayloadUserChangedStatus{
			JID:       jid,
			LID:       lid,
			Status:    evt.Status,
			Timestamp: evt.Timestamp,
		},
		&inst.ID,
	))
}

func (g *WhatsmeowGateway) emitUserPresence(inst *instance.Instance, evt *meowEvents.Presence) {
	lastSeen := &evt.LastSeen
	online := !evt.Unavailable
	jid, lid := GetJIDAndLID(evt.From)

	if lastSeen.IsZero() {
		lastSeen = nil
	}

	g.eventbus.Publish(events.New(
		user.EventChangedPresence,
		user.PayloadUserPresence{
			JID:      jid,
			LID:      lid,
			Online:   online,
			LastSeen: lastSeen,
		},
		&inst.ID,
	))
}

func (g *WhatsmeowGateway) emitUserNewMessage(inst *instance.Instance, evt *meowEvents.Message) {
	m := whatsmeowMessageToDomainMessage(evt)
	m.InstanceID = &inst.ID

	var eventName events.EventName

	switch m.Content.Kind() {
	case message.MessageKindText:
		eventName = user.EventNewTextMessage
	case message.MessageKindImage:
		eventName = user.EventNewImageMessage
	case message.MessageKindVideo:
		eventName = user.EventNewVideoMessage
	case message.MessageKindAudio:
		eventName = user.EventNewAudioMessage
	case message.MessageKindVoice:
		eventName = user.EventNewVoiceMessage
	case message.MessageKindDocument:
		eventName = user.EventNewDocumentMessage
	case message.MessageKindReaction:
		eventName = message.EventMessageReactionNew
		if reaction, ok := m.Content.(message.ReactionContent); ok {
			if reaction.Removed {
				eventName = message.EventMessageReactionRemoved
			}
		}
	default:
		fmt.Println("Unsupported message type:", m.Content.Kind())
		// Unsupported message type
		return
	}

	g.eventbus.Publish(events.New(
		eventName,
		message.PayloadNewMessage{
			Chat:    evt.Info.Chat.String(),
			Sender:  getSenderFromMessage(evt),
			Message: *m,
		},
		&inst.ID,
	))
}

// #region Newsletter Event Handlers
func (g *WhatsmeowGateway) emitNewsletterChangedPhoto(inst *instance.Instance, evt *meowEvents.Picture) {
	jid, lid := GetJIDAndLID(evt.Author)
	g.eventbus.Publish(events.New(
		newsletter.EventChangedPhoto,
		newsletter.PayloadNewsletterChangedPhoto{
			ID: evt.JID.String(),
			Changer: newsletter.Changer{
				JID: jid,
				LID: lid,
			},
			Photo:     evt.PictureID,
			Removed:   evt.Remove,
			Timestamp: evt.Timestamp,
		},
		&inst.ID,
	))
}

func (g *WhatsmeowGateway) emitNewsletterNewPost(inst *instance.Instance, evt *meowEvents.Message) {
	m := whatsmeowMessageToDomainMessage(evt)
	m.InstanceID = &inst.ID

	var eventName events.EventName

	switch m.Content.Kind() {
	case message.MessageKindText:
		eventName = newsletter.EventNewTextMessage
	case message.MessageKindImage:
		eventName = newsletter.EventNewImageMessage
	case message.MessageKindVideo:
		eventName = newsletter.EventNewVideoMessage
	case message.MessageKindAudio:
		eventName = newsletter.EventNewAudioMessage
	case message.MessageKindVoice:
		eventName = newsletter.EventNewVoiceMessage
	case message.MessageKindDocument:
		eventName = newsletter.EventNewDocumentMessage
	case message.MessageKindReaction:
		eventName = message.EventMessageReactionNew
		if reaction, ok := m.Content.(message.ReactionContent); ok {
			if reaction.Removed {
				eventName = message.EventMessageReactionRemoved
			}
		}
	default:
		fmt.Println("Unsupported message type:", m.Content.Kind())
		// Unsupported message type
		return
	}

	g.eventbus.Publish(events.New(
		eventName,
		message.PayloadNewMessage{
			Chat:    evt.Info.Chat.String(),
			Sender:  getSenderFromMessage(evt),
			Message: *m,
		},
		&inst.ID,
	))
}

// #region Group Event Handlers
func (g *WhatsmeowGateway) emitGroupChangedPhoto(inst *instance.Instance, evt *meowEvents.Picture) {
	jid, lid := GetJIDAndLID(evt.Author)
	g.eventbus.Publish(events.New(
		group.EventChangedPhoto,
		group.PayloadGroupChangedPhoto{
			JID: evt.JID.String(),
			Changer: group.Changer{
				JID: jid,
				LID: lid,
			},
			Photo:     evt.PictureID,
			Removed:   evt.Remove,
			Timestamp: evt.Timestamp,
		},
		&inst.ID,
	))
}

func (g *WhatsmeowGateway) emitGroupChangedInfo(inst *instance.Instance, evt *meowEvents.GroupInfo) {
	if len(evt.Join) > 0 {
		participants := make([]string, len(evt.Join))
		for i, p := range evt.Join {
			participants[i] = p.String()
		}

		g.eventbus.Publish(events.New(
			group.EventParticipantsJoined,
			group.PayloadGroupParticipantsJoined{
				JID:          evt.JID.String(),
				Participants: participants,
				Timestamp:    evt.Timestamp,
			},
			&inst.ID,
		))
	}

	if len(evt.Leave) > 0 {
		participants := make([]string, len(evt.Leave))
		for i, p := range evt.Leave {
			participants[i] = p.String()
		}

		g.eventbus.Publish(events.New(
			group.EventParticipantsLeft,
			group.PayloadGroupParticipantsLeft{
				JID:          evt.JID.String(),
				Participants: participants,
				Timestamp:    evt.Timestamp,
			},
			&inst.ID,
		))
	}

	if len(evt.Promote) > 0 {
		participants := make([]string, len(evt.Promote))
		for i, p := range evt.Promote {
			participants[i] = p.String()
		}

		jid, lid := "", ""
		if evt.Sender != nil {
			jid, lid = GetJIDAndLID(*evt.Sender)
		}

		g.eventbus.Publish(events.New(
			group.EventParticipantsPromoted,
			group.PayloadGroupParticipantsPromoted{
				JID:          evt.JID.String(),
				Participants: participants,
				Promoter:     group.Changer{JID: jid, LID: lid},
			},
			&inst.ID,
		))
	}

	if len(evt.Demote) > 0 {
		participants := make([]string, len(evt.Demote))
		for i, p := range evt.Demote {
			participants[i] = p.String()
		}

		jid, lid := "", ""
		if evt.Sender != nil {
			jid, lid = GetJIDAndLID(*evt.Sender)
		}

		g.eventbus.Publish(events.New(
			group.EventParticipantsPromoted,
			group.PayloadGroupParticipantsPromoted{
				JID:          evt.JID.String(),
				Participants: participants,
				Promoter:     group.Changer{JID: jid, LID: lid},
				Timestamp:    evt.Timestamp,
			},
			&inst.ID,
		))

		g.eventbus.Publish(events.New(
			group.EventParticipantsDemoted,
			group.PayloadGroupParticipantsDemoted{
				JID:          evt.JID.String(),
				Participants: participants,
				Demoter:      group.Changer{JID: jid, LID: lid},
				Timestamp:    evt.Timestamp,
			},
			&inst.ID,
		))
	}

	if evt.Name != nil {
		jid, lid := GetJIDAndLID(evt.Name.NameSetBy)
		if jid == "" {
			jid = evt.Name.NameSetByPN.String()
		}

		g.eventbus.Publish(events.New(
			group.EventChangedName,
			group.PayloadGroupNameChanged{
				JID:       evt.JID.String(),
				Name:      evt.Name.Name,
				Changer:   group.Changer{JID: jid, LID: lid},
				Timestamp: evt.Timestamp,
			},
			&inst.ID,
		))
	}

	if evt.Topic != nil {
		jid, lid := GetJIDAndLID(evt.Topic.TopicSetBy)
		if jid == "" {
			jid = evt.Topic.TopicSetByPN.String()
		}

		g.eventbus.Publish(events.New(
			group.EventChangedDescription,
			group.PayloadGroupDescriptionChanged{
				JID:         evt.JID.String(),
				Description: evt.Topic.Topic,
				Changer:     group.Changer{JID: jid, LID: lid},
				Deleted:     evt.Topic.TopicDeleted,
				Timestamp:   evt.Timestamp,
			},
			&inst.ID,
		))
	}

	if evt.Locked != nil {
		jid, lid := "", ""
		if evt.Sender != nil {
			jid, lid = GetJIDAndLID(*evt.Sender)
		}

		g.eventbus.Publish(events.New(
			group.EventChangedLocked,
			group.PayloadGroupChangedPermission{
				JID:       evt.JID.String(),
				Enabled:   evt.Locked.IsLocked,
				Changer:   group.Changer{JID: jid, LID: lid},
				Timestamp: evt.Timestamp,
			},
			&inst.ID,
		))
	}

	if evt.Announce != nil {
		jid, lid := "", ""
		if evt.Sender != nil {
			jid, lid = GetJIDAndLID(*evt.Sender)
		}

		g.eventbus.Publish(events.New(
			group.EventChangedAnnounce,
			group.PayloadGroupChangedPermission{
				JID:       evt.JID.String(),
				Enabled:   evt.Announce.IsAnnounce,
				Changer:   group.Changer{JID: jid, LID: lid},
				Timestamp: evt.Timestamp,
			},
			&inst.ID,
		))
	}

	if evt.MembershipApprovalMode != nil {
		jid, lid := "", ""
		if evt.Sender != nil {
			jid, lid = GetJIDAndLID(*evt.Sender)
		}
		g.eventbus.Publish(events.New(
			group.EventChangedApproval,
			group.PayloadGroupChangedPermission{
				JID:       evt.JID.String(),
				Enabled:   evt.MembershipApprovalMode.IsJoinApprovalRequired,
				Changer:   group.Changer{JID: jid, LID: lid},
				Timestamp: evt.Timestamp,
			},
			&inst.ID,
		))
	}

	if evt.Ephemeral != nil {
		jid, lid := "", ""
		if evt.Sender != nil {
			jid, lid = GetJIDAndLID(*evt.Sender)
		}

		g.eventbus.Publish(events.New(
			group.EventChangedExpiration,
			group.PayloadGroupExpirationChanged{
				JID:        evt.JID.String(),
				Enabled:    evt.Ephemeral.IsEphemeral,
				Expiration: evt.Ephemeral.DisappearingTimer,
				Duration:   group.NewGroupMessageDurationFromSeconds(evt.Ephemeral.DisappearingTimer),
				Changer:    group.Changer{JID: jid, LID: lid},
				Timestamp:  evt.Timestamp,
			},
			&inst.ID,
		))
	}
}

func (g *WhatsmeowGateway) emitGroupNewMessage(inst *instance.Instance, evt *meowEvents.Message) {
	m := whatsmeowMessageToDomainMessage(evt)
	m.InstanceID = &inst.ID

	var eventName events.EventName

	switch m.Content.Kind() {
	case message.MessageKindText:
		eventName = group.EventNewTextMessage
	case message.MessageKindImage:
		eventName = group.EventNewImageMessage
	case message.MessageKindVideo:
		eventName = group.EventNewVideoMessage
	case message.MessageKindAudio:
		eventName = group.EventNewAudioMessage
	case message.MessageKindVoice:
		eventName = group.EventNewVoiceMessage
	case message.MessageKindDocument:
		eventName = group.EventNewDocumentMessage
	case message.MessageKindReaction:
		eventName = message.EventMessageReactionNew
		if reaction, ok := m.Content.(message.ReactionContent); ok {
			if reaction.Removed {
				eventName = message.EventMessageReactionRemoved
			}
		}
	default:
		fmt.Println("Unsupported message type:", m.Content.Kind())
		return
	}

	g.eventbus.Publish(events.New(
		eventName,
		message.PayloadNewMessage{
			Chat:    evt.Info.Chat.String(),
			Sender:  getSenderFromMessage(evt),
			Message: *m,
		},
		&inst.ID,
	))
}

// #region Community Event Handlers
func (g *WhatsmeowGateway) emitCommunityChangedPhoto(inst *instance.Instance, evt *meowEvents.Picture) {
	jid, lid := GetJIDAndLID(evt.Author)
	g.eventbus.Publish(events.New(
		community.EventChangedPhoto,
		community.PayloadCommunityChangedPhoto{
			JID: evt.JID.String(),
			Changer: community.Changer{
				JID: jid,
				LID: lid,
			},
			Photo:     evt.PictureID,
			Removed:   evt.Remove,
			Timestamp: evt.Timestamp,
		},
		&inst.ID,
	))
}

func (g *WhatsmeowGateway) emitCommunityAnnouncement(inst *instance.Instance, evt *meowEvents.Message) {
	m := whatsmeowMessageToDomainMessage(evt)
	m.InstanceID = &inst.ID

	var eventName events.EventName

	switch m.Content.Kind() {
	case message.MessageKindText:
		eventName = community.EventNewTextMessage
	case message.MessageKindImage:
		eventName = community.EventNewImageMessage
	case message.MessageKindVideo:
		eventName = community.EventNewVideoMessage
	case message.MessageKindAudio:
		eventName = community.EventNewAudioMessage
	case message.MessageKindVoice:
		eventName = community.EventNewVoiceMessage
	case message.MessageKindDocument:
		eventName = community.EventNewDocumentMessage
	case message.MessageKindReaction:
		eventName = message.EventMessageReactionNew
		if reaction, ok := m.Content.(message.ReactionContent); ok {
			if reaction.Removed {
				eventName = message.EventMessageReactionRemoved
			}
		}
	default:
		fmt.Println("Unsupported message type:", m.Content.Kind())
		// Unsupported message type
		return
	}

	g.eventbus.Publish(events.New(
		eventName,
		message.PayloadNewMessage{
			Chat:    evt.Info.Chat.String(),
			Sender:  getSenderFromMessage(evt),
			Message: *m,
		},
		&inst.ID,
	))
}

// #region Chat Event Handlers
func (g *WhatsmeowGateway) emitChatPresence(inst *instance.Instance, evt *meowEvents.ChatPresence) {
	presence := chat.ChatPresencePaused

	if evt.Media == types.ChatPresenceMediaText && evt.State == types.ChatPresenceComposing {
		presence = chat.ChatPresenceTyping
	}

	if evt.Media == types.ChatPresenceMediaAudio && evt.State == types.ChatPresenceComposing {
		presence = chat.ChatPresenceRecording
	}

	jid, lid := GetJIDAndLID(evt.Sender)

	g.eventbus.Publish(events.New(
		chat.ChatChangedPresence,
		chat.PayloadChatChangedPresence{
			Chat:     evt.Chat.String(),
			Sender:   chat.Sender{JID: jid, LID: lid},
			Presence: presence,
		},
		&inst.ID,
	))
}

func (g *WhatsmeowGateway) emitChatRead(inst *instance.Instance, evt *meowEvents.MarkChatAsRead) {
	g.eventbus.Publish(events.New(
		chat.ChatRead,
		chat.PayloadChatStateRead{
			Chat:      evt.JID.String(),
			Timestamp: evt.Timestamp,
		},
		&inst.ID,
	))
}

func (g *WhatsmeowGateway) emitChatCleared(inst *instance.Instance, evt *meowEvents.ClearChat) {
	g.eventbus.Publish(events.New(
		chat.ChatCleared,
		chat.PayloadChatStateCleared{
			Chat:      evt.JID.String(),
			Timestamp: evt.Timestamp,
		},
		&inst.ID,
	))
}

func (g *WhatsmeowGateway) emitChatDeleted(inst *instance.Instance, evt *meowEvents.DeleteChat) {
	g.eventbus.Publish(events.New(
		chat.ChatDeleted,
		chat.PayloadChatStateDeleted{
			Chat:      evt.JID.String(),
			Timestamp: evt.Timestamp,
		},
		&inst.ID,
	))
}

func (g *WhatsmeowGateway) emitChatChangedPin(inst *instance.Instance, evt *meowEvents.Pin) {
	pinned := false
	if evt.Action.Pinned != nil {
		pinned = *evt.Action.Pinned
	}

	g.eventbus.Publish(events.New(
		chat.ChatChangedPin,
		chat.PayloadChatChangedPin{
			Chat:      evt.JID.String(),
			Pinned:    pinned,
			Timestamp: evt.Timestamp,
		},
		&inst.ID,
	))
}

func (g *WhatsmeowGateway) emitChatChangedMute(inst *instance.Instance, evt *meowEvents.Mute) {
	muted := false
	if evt.Action.Muted != nil {
		muted = *evt.Action.Muted
	}

	g.eventbus.Publish(events.New(
		chat.ChatChangedMute,
		chat.PayloadChatChangedMute{
			Chat:      evt.JID.String(),
			Muted:     muted,
			Timestamp: evt.Timestamp,
		},
		&inst.ID,
	))
}

func (g *WhatsmeowGateway) emitChatChangedArchive(inst *instance.Instance, evt *meowEvents.Archive) {
	archived := false
	if evt.Action.Archived != nil {
		archived = *evt.Action.Archived
	}

	g.eventbus.Publish(events.New(
		chat.ChatChangedArchive,
		chat.PayloadChatChangedArchive{
			Chat:      evt.JID.String(),
			Archived:  archived,
			Timestamp: evt.Timestamp,
		},
		&inst.ID,
	))
}

// #region Blocklist Event Emitters
func (g *WhatsmeowGateway) emitBlocklistChanged(inst *instance.Instance, evt *meowEvents.Blocklist) {
	changes := make([]blocklist.BlocklistChange, len(evt.Changes))
	for i, c := range evt.Changes {
		blocked := true
		if c.Action == meowEvents.BlocklistChangeActionUnblock {
			blocked = false
		}

		jid, lid := GetJIDAndLID(c.JID)
		changes[i] = blocklist.BlocklistChange{
			JID:     jid,
			LID:     lid,
			Blocked: blocked,
		}
	}
	g.eventbus.Publish(events.New(
		blocklist.EventChanged,
		blocklist.PayloadChanged{
			Changes: changes,
		},
		&inst.ID,
	))
}

func whatsmeowMessageToDomainMessage(msg *meowEvents.Message) *message.Message {
	var expiration *uint32
	var content message.Content

	if msg.Message.GetExtendedTextMessage() != nil {
		raw := msg.Message.GetExtendedTextMessage()

		mentioned := raw.GetContextInfo().GetMentionedJID()

		content = message.NewTextContent(
			raw.GetText(),
			&mentioned,
		)
	}

	if msg.Message.GetConversation() != "" {
		content = message.NewTextContent(
			msg.Message.GetConversation(),
			&[]string{},
		)
	}

	if msg.Message.GetImageMessage() != nil {
		raw := msg.Message.GetImageMessage()

		caption := raw.GetCaption()
		viewOnce := raw.GetViewOnce()
		mentioned := raw.GetContextInfo().GetMentionedJID()

		expiration = raw.GetContextInfo().Expiration

		width := raw.GetWidth()
		height := raw.GetHeight()
		imageFile := file.ImageFile{
			File: file.File{
				URL:       raw.GetURL(),
				Path:      raw.GetDirectPath(),
				Mime:      raw.GetMimetype(),
				Size:      raw.GetFileLength(),
				Sha256:    fmt.Sprintf("%x", raw.GetFileSHA256()),
				Sha256Enc: fmt.Sprintf("%x", raw.GetFileEncSHA256()),
				MediaKey:  fmt.Sprintf("%x", raw.GetMediaKey()),
				Extension: file.DetectExtension(raw.GetMimetype()),
			},
			Width:  &width,
			Height: &height,
		}

		thumbnail := new(string)
		if raw.GetJPEGThumbnail() != nil {
			*thumbnail = base64.StdEncoding.EncodeToString(raw.GetJPEGThumbnail())
		}

		content = message.NewImageContent(
			&imageFile,
			thumbnail,
			&caption,
			&mentioned,
			&viewOnce,
		)
	}

	if msg.Message.GetVideoMessage() != nil {
		raw := msg.Message.GetVideoMessage()

		caption := raw.GetCaption()
		viewOnce := raw.GetViewOnce()
		mentioned := raw.GetContextInfo().GetMentionedJID()
		expiration = raw.GetContextInfo().Expiration

		width := raw.GetWidth()
		height := raw.GetHeight()
		duration := raw.GetSeconds()
		videoFile := file.VideoFile{
			File: file.File{
				URL:       raw.GetURL(),
				Path:      raw.GetDirectPath(),
				Mime:      raw.GetMimetype(),
				Size:      raw.GetFileLength(),
				Sha256:    fmt.Sprintf("%x", raw.GetFileSHA256()),
				Sha256Enc: fmt.Sprintf("%x", raw.GetFileEncSHA256()),
				MediaKey:  fmt.Sprintf("%x", raw.GetMediaKey()),
				Extension: file.DetectExtension(raw.GetMimetype()),
			},
			Width:    &width,
			Height:   &height,
			Duration: &duration,
		}

		thumbnail := new(string)
		if raw.GetJPEGThumbnail() != nil {
			*thumbnail = base64.StdEncoding.EncodeToString(raw.GetJPEGThumbnail())
		}

		content = message.NewVideoContent(
			videoFile,
			thumbnail,
			&caption,
			&mentioned,
			&viewOnce,
		)
	}

	if msg.Message.AudioMessage != nil {
		raw := msg.Message.GetAudioMessage()

		duration := raw.GetSeconds()
		viewOnce := raw.GetViewOnce()
		expiration = raw.GetContextInfo().Expiration

		f := file.File{
			URL:       raw.GetURL(),
			Path:      raw.GetDirectPath(),
			Mime:      raw.GetMimetype(),
			Size:      raw.GetFileLength(),
			Sha256:    hex.EncodeToString(raw.GetFileSHA256()),
			Sha256Enc: hex.EncodeToString(raw.GetFileEncSHA256()),
			MediaKey:  hex.EncodeToString(raw.GetMediaKey()),
			Extension: file.DetectExtension(raw.GetMimetype()),
		}

		if msg.Message.AudioMessage.PTT != nil && *msg.Message.AudioMessage.PTT {
			content = message.NewVoiceContent(file.VoiceFile{File: f, Duration: &duration}, &viewOnce)
		} else {
			content = message.NewAudioContent(file.AudioFile{File: f, Duration: &duration})
		}
	}

	if msg.Message.GetDocumentMessage() != nil {
		raw := msg.Message.GetDocumentMessage()

		caption := raw.GetCaption()
		expiration = raw.GetContextInfo().Expiration

		pages := raw.GetPageCount()
		docFile := file.File{
			Name:      raw.GetFileName(),
			URL:       raw.GetURL(),
			Path:      raw.GetDirectPath(),
			Mime:      raw.GetMimetype(),
			Size:      raw.GetFileLength(),
			Sha256:    hex.EncodeToString(raw.GetFileSHA256()),
			Sha256Enc: hex.EncodeToString(raw.GetFileEncSHA256()),
			MediaKey:  hex.EncodeToString(raw.GetMediaKey()),
			Extension: file.DetectExtension(raw.GetMimetype()),
			Pages:     &pages,
		}

		thumbnail := new(string)
		if raw.GetJPEGThumbnail() != nil {
			*thumbnail = base64.StdEncoding.EncodeToString(raw.GetJPEGThumbnail())
		}

		content = message.NewDocumentContent(
			docFile,
			thumbnail,
			&caption,
			&raw.GetContextInfo().MentionedJID,
		)
	}

	if msg.Message.GetReactionMessage() != nil {
		raw := msg.Message.GetReactionMessage()
		content = message.NewReactionContent(
			raw.GetText(),
			*raw.GetKey().ID,
		)
	}

	message := message.NewMessage(
		&msg.Info.ID,
		msg.Info.Sender.String(),
		msg.Info.Chat.String(),
		content,
		nil,
		expiration,
		msg.Info.IsFromMe,
	)

	return message
}

func getSenderFromMessage(msg *meowEvents.Message) message.Sender {
	phone := ""
	jid, lid := GetJIDAndLID(msg.Info.Sender)
	if jid == "" {
		jid = msg.Info.SenderAlt.String()
		phone = msg.Info.SenderAlt.User
	} else {
		phone = msg.Info.Sender.User
	}

	if lid == "" {
		lid = msg.Info.SenderAlt.String()
	}

	return message.Sender{
		JID:   jid,
		LID:   lid,
		Phone: phone,
		Name:  msg.Info.PushName,
	}
}
