package service

import (
	"context"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/app/whatsapp"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/group"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/utils"
)

type GroupService struct {
	whatsapp    whatsapp.WhatsAppGateway
	eventbus    events.EventBus
	fileService *FileService
}

func NewGroupService(whatsapp whatsapp.WhatsAppGateway, eventbus events.EventBus, fileService *FileService) *GroupService {
	return &GroupService{
		whatsapp:    whatsapp,
		eventbus:    eventbus,
		fileService: fileService,
	}
}

func (s *GroupService) GetGroup(ctx context.Context, inst *instance.Instance, inp input.GetGroup) (*group.Group, *app.AppError) {
	l := app.GetGroupServiceLogger()
	l.Debug("Getting group", "instance", inst.ID, "group", inp.JID)

	if inp.JID == "" {
		return nil, app.TranslateError("group service", group.ErrNotFound)
	}

	g, err := s.whatsapp.GetGroup(ctx, inst, inp.JID, inp.WithParticipants)
	if err != nil {
		return nil, app.TranslateError("group service", err)
	}

	return g, nil
}

func (s *GroupService) GetGroups(ctx context.Context, inst *instance.Instance, inp input.GetGroups) ([]*group.Group, *app.AppError) {
	l := app.GetGroupServiceLogger()
	l.Debug("Getting groups", "instance", inst.ID)

	groups, err := s.whatsapp.GetGroups(ctx, inst, inp.WithParticipants)
	if err != nil {
		return []*group.Group{}, app.TranslateError("group service", err)
	}

	return groups, nil
}

func (s *GroupService) JoinGroup(ctx context.Context, inst *instance.Instance, inp input.JoinGroup) (*group.Group, *app.AppError) {
	l := app.GetGroupServiceLogger()
	l.Debug("Joining group", "instance", inst.ID, "invite", inp.Invite)

	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("group service", err)
	}

	group, err := s.whatsapp.JoinGroup(ctx, inst, inp.Invite)
	if err != nil {
		return nil, app.TranslateError("group service", err)
	}

	l.Info("Joined group successfully", "instance", inst.ID, "group", group.JID)

	return group, nil
}

func (s *GroupService) LeaveGroup(ctx context.Context, inst *instance.Instance, inp input.LeaveGroup) (*group.Group, *app.AppError) {
	l := app.GetGroupServiceLogger()
	l.Debug("Leaving group", "instance", inst.ID, "group", inp.JID)

	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("group service", err)
	}

	group, err := s.whatsapp.LeaveGroup(ctx, inst, inp.JID)
	if err != nil {
		return nil, app.TranslateError("group service", err)
	}

	l.Info("Left group successfully", "instance", inst.ID, "groupJID", inp.JID)

	return group, nil
}

func (s *GroupService) GetGroupInviteLink(ctx context.Context, inst *instance.Instance, inp input.GetGroupInvite) (string, *app.AppError) {
	l := app.GetGroupServiceLogger()
	l.Debug("Getting group invite link", "instance", inst.ID, "group", inp.JID, "refresh", inp.Refresh)

	if err := inp.Validate(); err != nil {
		return "", app.TranslateError("group service", err)
	}

	link, err := s.whatsapp.GetGroupInviteLink(ctx, inst, inp.JID, inp.Refresh)
	if err != nil {
		return "", app.TranslateError("group service", err)
	}

	return link, nil
}

func (s *GroupService) UpdateGroupName(ctx context.Context, inst *instance.Instance, inp input.UpdateGroupName) (*group.Group, *app.AppError) {
	l := app.GetGroupServiceLogger()
	l.Debug("Updating group name", "instance", inst.ID, "group", inp.JID, "name", inp.Name)

	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("group service", err)
	}

	g, err := s.whatsapp.GetGroup(ctx, inst, inp.JID, utils.BoolPtr(false))
	if err != nil {
		return nil, app.TranslateError("group service", err)
	}

	err = s.whatsapp.GroupUpdateName(ctx, inst, inp.JID, inp.Name)
	if err != nil {
		return nil, app.TranslateError("group service", err)
	}

	g.Name = inp.Name

	l.Info("Updated group name successfully", "instance", inst.ID, "groupJID", inp.JID, "name", inp.Name)

	return g, nil
}

func (s *GroupService) UpdateGroupDescription(ctx context.Context, inst *instance.Instance, inp input.UpdateGroupDescription) (*group.Group, *app.AppError) {
	l := app.GetGroupServiceLogger()
	l.Debug("Updating group description", "instance", inst.ID, "groupJID", inp.JID, "description", inp.Description)

	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("group service", err)
	}

	g, err := s.whatsapp.GetGroup(ctx, inst, inp.JID, utils.BoolPtr(false))
	if err != nil {
		return nil, app.TranslateError("group service", err)
	}

	err = s.whatsapp.GroupUpdateDescription(ctx, inst, inp.JID, inp.Description)
	if err != nil {
		return nil, app.TranslateError("group service", err)
	}

	g.Description = inp.Description

	l.Info("Updated group description successfully", "instance", inst.ID, "groupJID", inp.JID, "description", inp.Description)

	return g, nil
}

func (s *GroupService) UpdateGroupSetting(ctx context.Context, inst *instance.Instance, inp input.UpdateGroupSetting) (*group.Group, *app.AppError) {
	l := app.GetGroupServiceLogger()
	l.Debug("Updating group setting", "instance", inst.ID, "group", inp.JID, "setting", inp.Setting, "policy", inp.Policy)

	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("group service", err)
	}

	g, err := s.whatsapp.GetGroup(ctx, inst, inp.JID, utils.BoolPtr(false))
	if err != nil {
		return nil, app.TranslateError("group service", err)
	}

	err = s.whatsapp.GroupUpdateSetting(ctx, inst, inp.JID, inp.Setting, inp.Policy)
	if err != nil {
		return nil, app.TranslateError("group service", err)
	}

	l.Info("Updated group setting successfully", "instance", inst.ID, "group", inp.JID, "setting", inp.Setting, "policy", inp.Policy)

	return g, nil
}

func (s *GroupService) UpdateGroupMessageDuration(ctx context.Context, inst *instance.Instance, inp input.UpdateGroupMessageDuration) (*group.Group, *app.AppError) {
	l := app.GetGroupServiceLogger()

	l.Debug("Updating group message duration", "instance", inst.ID, "group", inp.JID, "duration", inp.Duration)

	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("group service", err)
	}

	g, err := s.whatsapp.GetGroup(ctx, inst, inp.JID, utils.BoolPtr(false))
	if err != nil {
		return nil, app.TranslateError("group service", err)
	}

	err = s.whatsapp.GroupUpdateMessageDuration(ctx, inst, inp.JID, inp.Duration)
	if err != nil {
		return nil, app.TranslateError("group service", err)
	}

	l.Info("Updated group message duration successfully", "instance", inst.ID, "group", inp.JID, "duration", inp.Duration)

	return g, nil
}

func (s *GroupService) UpdateGroupParticipants(ctx context.Context, inst *instance.Instance, inp input.UpdateGroupParticipants) ([]*group.GroupParticipant, *app.AppError) {
	l := app.GetGroupServiceLogger()
	l.Debug("Updating group participants", "instance", inst.ID, "group", inp.JID, "participants", inp.Participants, "action", inp.Action)

	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("group service", err)
	}

	updatedParticipants, err := s.whatsapp.GroupUpdateParticipants(ctx, inst, inp.JID, inp.Participants, inp.Action)
	if err != nil {
		return nil, app.TranslateError("group service", err)
	}

	l.Info("Updated group participants successfully", "instance", inst.ID, "group", inp.JID, "participants", inp.Participants, "action", inp.Action)

	return updatedParticipants, nil
}

func (s *GroupService) UpdateGroupPhoto(ctx context.Context, inst *instance.Instance, inp input.UpdateGroupPhoto) (string, *app.AppError) {
	l := app.GetGroupServiceLogger()
	l.Debug("Updating group photo", "instance", inst.ID, "group", inp.JID)

	if err := inp.Validate(); err != nil {
		return "", app.TranslateError("group service", err)
	}

	_, data, err := s.fileService.GetFrom(ctx, inp.Photo)
	if err != nil {
		return "", app.TranslateError("group service", err)
	}

	groupPhoto, err := group.NewGroupPhoto(*data)
	if err != nil {
		return "", app.TranslateError("group service", err)
	}

	newPhotoID, err := s.whatsapp.GroupUpdatePhoto(ctx, inst, inp.JID, groupPhoto)
	if err != nil {
		return "", app.TranslateError("group service", err)
	}

	l.Info("Updated group photo successfully", "instance", inst.ID, "group", inp.JID)

	return newPhotoID, nil
}

func (s *GroupService) RemoveGroupPhoto(ctx context.Context, inst *instance.Instance, inp input.RemoveGroupPhoto) *app.AppError {
	l := app.GetGroupServiceLogger()
	l.Debug("Removing group photo", "instance", inst.ID, "group", inp.JID)

	if err := inp.Validate(); err != nil {
		return app.TranslateError("group service", err)
	}

	err := s.whatsapp.GroupRemovePhoto(ctx, inst, inp.JID)
	if err != nil {
		return app.TranslateError("group service", err)
	}

	l.Info("Removed group photo successfully", "instance", inst.ID, "group", inp.JID)

	return nil
}

func (s *GroupService) GetGroupPhotoURL(ctx context.Context, inst *instance.Instance, inp input.GetPhotoURL) (string, *app.AppError) {
	l := app.GetGroupServiceLogger()
	l.Debug("Getting group photo URL", "instance", inst.ID, "group", inp.JID)

	if err := inp.Validate(); err != nil {
		return "", app.TranslateError("group service", err)
	}

	g, err := s.whatsapp.GetGroup(ctx, inst, inp.JID, utils.BoolPtr(false))
	if err != nil {
		return "", app.TranslateError("group service", err)
	}

	photoURL, err := s.whatsapp.GetPictureURL(ctx, inst, inp.JID, inp.Preview, g.Type == group.GroupTypeCommunity)
	if err != nil {
		return "", app.TranslateError("group service", err)
	}

	return photoURL, nil
}

func (s *GroupService) CreateGroup(ctx context.Context, inst *instance.Instance, inp input.CreateGroup) (*group.Group, *app.AppError) {
	l := app.GetGroupServiceLogger()
	l.Debug("Creating group", "instance", inst.ID, "name", inp.Name)

	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("group service", err)
	}

	g, err := s.whatsapp.GroupCreate(ctx, inst, inp.Name, inp.Participants)
	if err != nil {
		return nil, app.TranslateError("group service", err)
	}

	l.Info("Created group successfully", "instance", inst.ID, "group", g.JID, "name", inp.Name)

	return g, nil
}
