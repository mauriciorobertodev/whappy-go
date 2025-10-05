package handler

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/app/service"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/group"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http/middleware"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http/requests"
	"github.com/mauriciorobertodev/whappy-go/internal/utils"
)

type GroupHandler struct {
	groupService *service.GroupService
	eventbus     events.EventBus
}

func NewGroupHandler(groupService *service.GroupService, eventbus events.EventBus) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
		eventbus:     eventbus,
	}
}

func (h *GroupHandler) RegisterRoutes(r fiber.Router, authMiddleware *middleware.AuthMiddleware, instMiddleware *middleware.InstanceMiddleware) {
	grp := r.Group("/groups", authMiddleware.Authenticate(), instMiddleware.AttachInstance(), instMiddleware.ConnectInstance())

	grp.Get("/", h.GetGroups)
	grp.Post("/", h.CreateGroup)
	grp.Get("/:group", h.GetGroup)
	grp.Delete("/:group", h.LeaveGroup)
	grp.Patch("/:group", h.UpdateGroupSettings)

	grp.Patch("/:group/name", h.UpdateGroupName)
	grp.Patch("/:group/description", h.UpdateGroupDescription)
	grp.Patch("/:group/disappearing", h.UpdateGroupMessageDuration)

	grp.Get("/:group/participants", h.GetParticipants)
	grp.Post("/:group/participants", h.AddParticipants)
	grp.Delete("/:group/participants", h.RemoveParticipants)

	grp.Post("/:group/admins", h.AddAdmins)
	grp.Delete("/:group/admins", h.RemoveAdmins)

	grp.Post("/join", h.JoinGroup)

	grp.Get("/:group/invite", h.GetGroupInviteLink)
	grp.Delete("/:group/invite", h.RevokeGroupInviteLink)

	grp.Get("/:group/photo", h.GetGroupPhoto)
	grp.Put("/:group/photo", h.UpdateGroupPhoto)
	grp.Delete("/:group/photo", h.RemoveGroupPhoto)
}

func (h *GroupHandler) GetGroup(c fiber.Ctx) error {
	inst := c.Locals("instance").(*instance.Instance)
	groupID := c.Params("group")
	withParticipants := c.Query("participants", "false") == "true"

	group, err := h.groupService.GetGroup(context.Background(), inst, input.GetGroup{
		JID:              groupID,
		WithParticipants: &withParticipants,
	})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to get group", err))
	}

	return c.JSON(http.NewSuccessResponse("Group retrieved successfully", fiber.Map{
		"group": group,
	}))
}

func (h *GroupHandler) GetGroups(c fiber.Ctx) error {
	inst := c.Locals("instance").(*instance.Instance)
	withParticipants := c.Query("participants", "false") == "true"

	groups, err := h.groupService.GetGroups(context.Background(), inst, input.GetGroups{
		WithParticipants: &withParticipants,
	})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to get groups", err))
	}

	return c.JSON(http.NewSuccessResponse("Groups retrieved successfully", fiber.Map{
		"groups": groups,
	}))
}

func (h *GroupHandler) JoinGroup(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)

	var req requests.JoinGroupRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	group, err := h.groupService.JoinGroup(ctx, inst, input.JoinGroup{Invite: req.Invite})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to join group", err))
	}

	return c.JSON(http.NewSuccessResponse("Joined group successfully", fiber.Map{
		"group": group,
	}))
}

func (h *GroupHandler) LeaveGroup(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)
	groupJID := c.Params("group")

	if groupJID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Group JID is required", nil))
	}

	group, err := h.groupService.LeaveGroup(ctx, inst, input.LeaveGroup{JID: groupJID})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to leave group", err))
	}

	return c.JSON(http.NewSuccessResponse("Left group successfully", fiber.Map{
		"group": group,
	}))
}

func (h *GroupHandler) GetGroupInviteLink(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)
	groupJID := c.Params("group")

	link, err := h.groupService.GetGroupInviteLink(ctx, inst, input.GetGroupInvite{JID: groupJID, Refresh: false})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to get group invite link", err))
	}

	return c.JSON(http.NewSuccessResponse("Group invite link retrieved successfully", fiber.Map{
		"link": link,
	}))
}

func (h *GroupHandler) RevokeGroupInviteLink(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)
	groupJID := c.Params("group")

	link, err := h.groupService.GetGroupInviteLink(ctx, inst, input.GetGroupInvite{JID: groupJID, Refresh: true})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to revoke group invite link", err))
	}

	return c.JSON(http.NewSuccessResponse("Group invite link revoked successfully", fiber.Map{
		"link": link,
	}))
}

func (h *GroupHandler) UpdateGroupSettings(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)
	groupJID := c.Params("group")

	var req requests.GroupUpdateSettingRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	if bag := req.Validate(); bag.HasErrors() {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewValidationErrorResponse(bag))
	}

	g, appErr := h.groupService.UpdateGroupSetting(ctx, inst, input.UpdateGroupSetting{
		JID:     groupJID,
		Setting: req.Setting,
		Policy:  req.Policy,
	})

	if appErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to update group settings", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Group settings updated successfully", fiber.Map{
		"group": g,
	}))
}

func (h *GroupHandler) UpdateGroupMessageDuration(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)
	groupJID := c.Params("group")

	var req requests.GroupUpdateMessageDurationRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	if bag := req.Validate(); bag.HasErrors() {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewValidationErrorResponse(bag))
	}

	g, appErr := h.groupService.UpdateGroupMessageDuration(ctx, inst, input.UpdateGroupMessageDuration{
		JID:      groupJID,
		Duration: req.Duration,
	})

	if appErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to update group message duration", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Group message duration updated successfully", fiber.Map{
		"group": g,
	}))
}

func (h *GroupHandler) AddParticipants(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)
	groupJID := c.Params("group")

	var req requests.UpdateParticipantsRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	if bag := req.Validate(); bag.HasErrors() {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewValidationErrorResponse(bag))
	}

	participants, appErr := h.groupService.UpdateGroupParticipants(ctx, inst, input.UpdateGroupParticipants{
		JID:          groupJID,
		Participants: req.Participants,
		Action:       group.ParticipantsActionAdd,
	})

	if appErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to update group participants", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Group participants updated successfully", fiber.Map{
		"participants": participants,
	}))
}

func (h *GroupHandler) RemoveParticipants(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)
	groupJID := c.Params("group")

	var req requests.UpdateParticipantsRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	if bag := req.Validate(); bag.HasErrors() {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewValidationErrorResponse(bag))
	}

	participants, appErr := h.groupService.UpdateGroupParticipants(ctx, inst, input.UpdateGroupParticipants{
		JID:          groupJID,
		Participants: req.Participants,
		Action:       group.ParticipantsActionRemove,
	})

	if appErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to update group participants", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Group participants updated successfully", fiber.Map{
		"participants": participants,
	}))
}

func (h *GroupHandler) GetParticipants(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)
	groupJID := c.Params("group")

	group, err := h.groupService.GetGroup(ctx, inst, input.GetGroup{
		JID:              groupJID,
		WithParticipants: utils.BoolPtr(true),
	})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to get group", err))
	}

	return c.JSON(http.NewSuccessResponse("Group participants retrieved successfully", fiber.Map{
		"participants": group.Participants,
		"count":        len(group.Participants),
	}))
}

func (h *GroupHandler) AddAdmins(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)
	groupJID := c.Params("group")

	var req requests.UpdateParticipantsRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	if bag := req.Validate(); bag.HasErrors() {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewValidationErrorResponse(bag))
	}

	participants, appErr := h.groupService.UpdateGroupParticipants(ctx, inst, input.UpdateGroupParticipants{
		JID:          groupJID,
		Participants: req.Participants,
		Action:       group.ParticipantsActionPromote,
	})

	if appErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to update group admins", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Group admins updated successfully", fiber.Map{
		"participants": participants,
	}))
}

func (h *GroupHandler) RemoveAdmins(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)
	groupJID := c.Params("group")

	var req requests.UpdateParticipantsRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	if bag := req.Validate(); bag.HasErrors() {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewValidationErrorResponse(bag))
	}

	participants, appErr := h.groupService.UpdateGroupParticipants(ctx, inst, input.UpdateGroupParticipants{
		JID:          groupJID,
		Participants: req.Participants,
		Action:       group.ParticipantsActionDemote,
	})

	if appErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to update group admins", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Group admins updated successfully", fiber.Map{
		"participants": participants,
	}))
}

func (h *GroupHandler) UpdateGroupPhoto(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)
	groupJID := c.Params("group")

	var req requests.GroupUpdatePhotoRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	if bag := req.Validate(); bag.HasErrors() {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewValidationErrorResponse(bag))
	}

	newPhotoID, appErr := h.groupService.UpdateGroupPhoto(ctx, inst, input.UpdateGroupPhoto{
		JID:   groupJID,
		Photo: req.Photo,
	})

	if appErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to update group photo", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Group photo updated successfully", fiber.Map{
		"photo": newPhotoID,
	}))
}

func (h *GroupHandler) RemoveGroupPhoto(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)
	groupJID := c.Params("group")

	appErr := h.groupService.RemoveGroupPhoto(ctx, inst, input.RemoveGroupPhoto{
		JID: groupJID,
	})

	if appErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to remove group photo", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Group photo removed successfully", nil))
}

func (h *GroupHandler) GetGroupPhoto(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)
	groupID := c.Params("group")
	preview := c.Query("preview", "false")

	photo, appErr := h.groupService.GetGroupPhotoURL(ctx, inst, input.GetPhotoURL{
		JID:     groupID,
		Preview: preview == "true",
	})

	if appErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to get group photo", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Group photo retrieved successfully", fiber.Map{
		"photo": photo,
	}))
}

func (h *GroupHandler) CreateGroup(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)

	var req requests.CreateGroupRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	if bag := req.Validate(); bag.HasErrors() {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewValidationErrorResponse(bag))
	}

	group, appErr := h.groupService.CreateGroup(ctx, inst, req.ToInput())
	if appErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to create group", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Group created successfully", fiber.Map{
		"group": group,
	}))
}

func (h *GroupHandler) UpdateGroupName(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)
	groupJID := c.Params("group")

	var req requests.UpdateGroupNameRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	if bag := req.Validate(); bag.HasErrors() {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewValidationErrorResponse(bag))
	}

	group, appErr := h.groupService.UpdateGroupName(ctx, inst, req.ToInput(groupJID))
	if appErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to update group name", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Group name updated successfully", fiber.Map{
		"group": group,
	}))
}

func (h *GroupHandler) UpdateGroupDescription(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)
	groupJID := c.Params("group")

	var req requests.UpdateGroupDescriptionRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	if bag := req.Validate(); bag.HasErrors() {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewValidationErrorResponse(bag))
	}

	group, appErr := h.groupService.UpdateGroupDescription(ctx, inst, req.ToInput(groupJID))
	if appErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to update group description", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Group description updated successfully", fiber.Map{
		"group": group,
	}))
}
