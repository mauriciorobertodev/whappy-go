package meow

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/group"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

func (g *WhatsmeowGateway) GetGroups(ctx context.Context, inst *instance.Instance, withParticipants *bool) ([]*group.Group, error) {
	l := app.GetWhatsappLogger()

	l.Info("Getting groups", "instance", inst.ID, "withParticipants", *withParticipants)

	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return nil, err
	}

	groupInfos, err := client.GetJoinedGroups(ctx)
	if err != nil {
		return nil, err
	}

	groups := make([]*group.Group, 0)
	for _, groupInfo := range groupInfos {
		group := groupInfoToGroup(groupInfo, inst, withParticipants)
		groups = append(groups, group)
	}

	return groups, nil
}

func (g *WhatsmeowGateway) GetGroup(ctx context.Context, inst *instance.Instance, groupJID string, withParticipants *bool) (*group.Group, error) {
	l := app.GetWhatsappLogger()

	l.Info("Getting group", "instance", inst.ID, "group", groupJID, "withParticipants", *withParticipants)

	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return nil, err
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return nil, err
	}

	groupInfo, err := client.GetGroupInfo(jid)
	if err != nil {
		if strings.Contains(err.Error(), "not a group") {
			return nil, group.ErrNotFound
		}
		return nil, err
	}

	group := groupInfoToGroup(groupInfo, inst, withParticipants)

	return group, nil
}

func (g *WhatsmeowGateway) JoinGroup(ctx context.Context, inst *instance.Instance, inviteCode string) (*group.Group, error) {
	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return nil, err
	}

	groupJID, err := client.JoinGroupWithLink(inviteCode)
	if err != nil {
		if errors.Is(err, whatsmeow.ErrInviteLinkInvalid) {
			return nil, group.ErrInviteInvalid
		}

		if errors.Is(err, whatsmeow.ErrInviteLinkRevoked) {
			return nil, group.ErrInviteRevoked
		}

		return nil, err
	}

	groupInfo, err := client.GetGroupInfo(groupJID)
	if err != nil {
		return nil, err
	}

	withParticipants := true
	return groupInfoToGroup(groupInfo, inst, &withParticipants), nil
}

func (w *WhatsmeowGateway) LeaveGroup(ctx context.Context, inst *instance.Instance, groupJID string) (*group.Group, error) {
	client, err := w.getOnlineClient(inst.ID)
	if err != nil {
		return nil, err
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return nil, err
	}

	withParticipants := false
	g, err := w.GetGroup(ctx, inst, groupJID, &withParticipants)
	if err != nil {
		return nil, err
	}

	if err := client.LeaveGroup(jid); err != nil {
		if err.Error() == "you're not participating in that group" {
			return nil, group.ErrNotMember
		}
		return nil, err
	}

	return g, nil
}

func (g *WhatsmeowGateway) GetGroupInviteLink(ctx context.Context, inst *instance.Instance, groupJID string, refresh bool) (string, error) {
	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return "", err
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return "", err
	}

	link, err := client.GetGroupInviteLink(jid, refresh)
	if err != nil {
		if errors.Is(err, whatsmeow.ErrGroupInviteLinkUnauthorized) {
			return "", group.ErrInviteUnauthorized
		}

		if errors.Is(err, whatsmeow.ErrGroupNotFound) {
			return "", group.ErrNotFound
		}

		if errors.Is(err, whatsmeow.ErrNotInGroup) {
			return "", group.ErrNotMember
		}

		return "", err
	}

	return link, nil
}

// Updates the group photo
func (g *WhatsmeowGateway) GroupUpdatePhoto(ctx context.Context, inst *instance.Instance, groupJID string, photo *group.GroupPhoto) (string, error) {
	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return "", err
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return "", err
	}

	newPhotoID, err := client.SetGroupPhoto(jid, photo.Data)
	if err != nil {
		if errors.Is(err, whatsmeow.ErrGroupNotFound) {
			return "", group.ErrNotFound
		}

		if errors.Is(err, whatsmeow.ErrNotInGroup) {
			return "", group.ErrNotMember
		}

		if errors.Is(err, whatsmeow.ErrInvalidImageFormat) {
			return "", group.ErrPhotoUnsupportedFormat
		}

		return "", err
	}

	return newPhotoID, nil
}

func (g *WhatsmeowGateway) GroupRemovePhoto(ctx context.Context, inst *instance.Instance, groupJID string) error {
	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return err
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return err
	}

	if _, err := client.SetGroupPhoto(jid, nil); err != nil {
		if errors.Is(err, whatsmeow.ErrGroupNotFound) {
			return group.ErrNotFound
		}

		if errors.Is(err, whatsmeow.ErrNotInGroup) {
			return group.ErrNotMember
		}

		return err
	}

	return nil
}

// Updates the group name
func (g *WhatsmeowGateway) GroupUpdateName(ctx context.Context, inst *instance.Instance, groupJID string, name string) error {
	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return err
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return err
	}

	if err := client.SetGroupName(jid, name); err != nil {
		if errors.Is(err, whatsmeow.ErrGroupNotFound) {
			return group.ErrNotFound
		}

		if errors.Is(err, whatsmeow.ErrNotInGroup) {
			return group.ErrNotMember
		}

		return err
	}

	return nil
}

// Updates the group description
func (g *WhatsmeowGateway) GroupUpdateDescription(ctx context.Context, inst *instance.Instance, groupJID string, description string) error {
	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return err
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return err
	}

	if err := client.SetGroupDescription(jid, description); err != nil {
		if errors.Is(err, whatsmeow.ErrGroupNotFound) {
			return group.ErrNotFound
		}

		if errors.Is(err, whatsmeow.ErrNotInGroup) {
			return group.ErrNotMember
		}

		return err
	}

	return nil
}

// Set the group topic, only available for communities
func (g *WhatsmeowGateway) GroupUpdateTopic(ctx context.Context, inst *instance.Instance, groupJID string, topic string) error {
	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return err
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return err
	}

	if err := client.SetGroupTopic(jid, "", "", topic); err != nil {
		if errors.Is(err, whatsmeow.ErrGroupNotFound) {
			return group.ErrNotFound
		}

		if errors.Is(err, whatsmeow.ErrNotInGroup) {
			return group.ErrNotMember
		}

		return err
	}

	return nil
}

// When enabled, only admins can update group info
func (g *WhatsmeowGateway) GroupUpdateSetting(ctx context.Context, inst *instance.Instance, groupJID string, setting group.GroupSettingName, policy group.GroupSettingPolicy) error {
	var err error
	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return err
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return err
	}

	enabled := false
	switch setting {
	case group.GroupSettingEditGroupInfo:
		if policy == group.GroupSettingPolicyAdmins {
			enabled = true
		}
		err = client.SetGroupLocked(jid, enabled)
	case group.GroupSettingSendMessages:
		if policy == group.GroupSettingPolicyAdmins {
			enabled = true
		}
		err = client.SetGroupAnnounce(jid, enabled)
	case group.GroupSettingAddParticipants:
		if policy == group.GroupSettingPolicyAdmins {
			err = client.SetGroupMemberAddMode(jid, types.GroupMemberAddModeAdmin)
		} else {
			err = client.SetGroupMemberAddMode(jid, types.GroupMemberAddModeAllMember)
		}
	case group.GroupSettingApproveParticipants:
		if policy == group.GroupSettingPolicyAdmins {
			enabled = true
		}
		err = client.SetGroupJoinApprovalMode(jid, enabled)
	}

	if err != nil {
		if errors.Is(err, whatsmeow.ErrGroupNotFound) {
			return group.ErrNotFound
		}

		if errors.Is(err, whatsmeow.ErrNotInGroup) {
			return group.ErrNotMember
		}

		return err
	}

	return nil
}

// Sets the ephemeral messages timer for the group
func (g *WhatsmeowGateway) GroupUpdateMessageDuration(ctx context.Context, inst *instance.Instance, groupJID string, duration group.GroupMessageDuration) error {
	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return err
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return err
	}

	if err := client.SetDisappearingTimer(jid, duration.ToDuration(), time.Time{}); err != nil {
		if errors.Is(err, whatsmeow.ErrGroupNotFound) {
			return group.ErrNotFound
		}

		if errors.Is(err, whatsmeow.ErrNotInGroup) {
			return group.ErrNotMember
		}

		return err
	}

	return nil
}

func (g *WhatsmeowGateway) GroupUpdateParticipants(ctx context.Context, inst *instance.Instance, groupJID string, participants []string, action group.ParticipantsAction) ([]*group.GroupParticipant, error) {
	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return nil, err
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return nil, err
	}

	jids := make([]types.JID, 0, len(participants))
	for _, p := range participants {
		parsedJID, err := types.ParseJID(p)
		if err != nil {
			return nil, err
		}
		jids = append(jids, parsedJID)
	}

	var updatedParticipants []types.GroupParticipant
	switch action {
	case group.ParticipantsActionAdd:
		updatedParticipants, err = client.UpdateGroupParticipants(jid, jids, whatsmeow.ParticipantChangeAdd)
	case group.ParticipantsActionRemove:
		updatedParticipants, err = client.UpdateGroupParticipants(jid, jids, whatsmeow.ParticipantChangeRemove)
	case group.ParticipantsActionPromote:
		updatedParticipants, err = client.UpdateGroupParticipants(jid, jids, whatsmeow.ParticipantChangePromote)
	case group.ParticipantsActionDemote:
		updatedParticipants, err = client.UpdateGroupParticipants(jid, jids, whatsmeow.ParticipantChangeDemote)
	default:
		return nil, group.ErrInvalidAction
	}

	if err != nil {
		if errors.Is(err, whatsmeow.ErrGroupNotFound) {
			return nil, group.ErrNotFound
		}

		if errors.Is(err, whatsmeow.ErrNotInGroup) {
			return nil, group.ErrNotMember
		}

		return nil, err
	}

	result := make([]*group.GroupParticipant, 0, len(updatedParticipants))
	for _, p := range updatedParticipants {
		result = append(result, participantInfoToParticipant(p, inst))
	}

	return result, nil
}

func (g *WhatsmeowGateway) GroupCreate(ctx context.Context, inst *instance.Instance, name string, participants []string) (*group.Group, error) {
	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return nil, err
	}

	jids := make([]types.JID, 0, len(participants))
	for _, p := range participants {
		if !strings.Contains(p, "@") {
			p = p + "@" + types.DefaultUserServer
		}

		parsedJID, err := types.ParseJID(p)
		if err != nil {
			return nil, err
		}
		jids = append(jids, parsedJID)
	}

	// TODO: can be improved to support settings
	req := whatsmeow.ReqCreateGroup{
		Name:         name,
		Participants: jids,
	}

	groupInfo, err := client.CreateGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	withParticipants := true
	return groupInfoToGroup(groupInfo, inst, &withParticipants), nil
}

func groupInfoToGroup(info *types.GroupInfo, inst *instance.Instance, withParticipants *bool) *group.Group {
	groupType := group.GroupTypeRegular
	if info.IsParent {
		groupType = group.GroupTypeCommunity
	} else if info.IsDefaultSubGroup && info.IsAnnounce {
		groupType = group.GroupTypeAnnouncement
	} else if !info.LinkedParentJID.IsEmpty() {
		groupType = group.GroupTypeSubGroup
	}

	messageExpiration := group.GroupMessageExpiration{
		Enabled:    info.GroupEphemeral.IsEphemeral,
		Expiration: info.GroupEphemeral.DisappearingTimer,
		Duration:   group.NewGroupMessageDurationFromSeconds(info.GroupEphemeral.DisappearingTimer),
	}

	jid, lid := GetJIDAndLID(info.OwnerJID)
	if jid == "" {
		jid = info.OwnerPN.String()
	}
	phone := info.OwnerPN.User

	g := &group.Group{
		JID: info.JID.String(),
		Owner: group.Owner{
			JID:   jid,
			LID:   lid,
			Phone: phone,
		},
		Type: groupType,

		Name:        info.Name,
		Description: info.Topic,

		Participation: group.GroupParticipationNone,

		Locked:     info.IsLocked,
		Announce:   info.IsAnnounce,
		Approval:   info.IsJoinApprovalRequired,
		Restricted: info.MemberAddMode == types.GroupMemberAddModeAdmin,
		Hidden:     info.IsIncognito,

		MessageExpiration: messageExpiration,

		ParentJID: info.LinkedParentJID.String(),

		AddressingMode: group.AddressingMode(info.AddressingMode),

		Participants: make([]*group.GroupParticipant, 0, len(info.Participants)),

		CreatorCountryCode: info.CreatorCountryCode,

		CreatedAt: info.GroupCreated.UTC(),
	}

	if withParticipants != nil && *withParticipants {
		for _, p := range info.Participants {
			g.Participants = append(g.Participants, participantInfoToParticipant(p, inst))
		}
	}

	return g
}

func participantInfoToParticipant(info types.GroupParticipant, inst *instance.Instance) *group.GroupParticipant {
	isMe := info.JID.String() == inst.JID || info.PhoneNumber.String() == inst.JID
	if !isMe && info.PhoneNumber.User != "" {
		isMe = info.PhoneNumber.User == inst.Phone
	}
	isOwner := info.JID.String() == inst.JID || info.PhoneNumber.String() == inst.JID
	if !isOwner && info.PhoneNumber.User != "" && inst.JID != "" {
		isOwner = info.PhoneNumber.User == inst.JID
	}
	isAdmin := info.IsAdmin || isOwner
	isSuperAdmin := info.IsSuperAdmin || isOwner

	return &group.GroupParticipant{
		JID:   info.PhoneNumber.String(),
		LID:   info.LID.String(),
		Phone: info.PhoneNumber.User,
		Name:  info.DisplayName,
		IsMe:  isMe,
		Participation: func() group.GroupParticipation {
			if isOwner {
				return group.GroupParticipationOwner
			} else if isSuperAdmin {
				return group.GroupParticipationSuperAdmin
			} else if isAdmin {
				return group.GroupParticipationAdmin
			}
			return group.GroupParticipationMember
		}(),
	}
}
