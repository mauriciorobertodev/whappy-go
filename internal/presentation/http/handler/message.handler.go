package handler

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/app/service"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http/middleware"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http/requests"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

func (h *MessageHandler) RegisterRoutes(r fiber.Router, authMiddleware *middleware.AuthMiddleware, instMiddleware *middleware.InstanceMiddleware) {
	msg := r.Group("/messages", authMiddleware.Authenticate(), instMiddleware.AttachInstance(), instMiddleware.ConnectInstance())
	msg.Get("/id", h.GetMessageIDs)
	msg.Post("/text", h.SendText)
	msg.Post("/image", h.SendImage)
	msg.Post("/video", h.SendVideo)
	msg.Post("/audio", h.SendAudio)
	msg.Post("/voice", h.SendVoice)
	msg.Post("/document", h.SendDocument)
	msg.Post("/reaction", h.SendReaction)
	msg.Post("/read", h.MarkMessagesAsRead)
}

func (h *MessageHandler) GetMessageIDs(c fiber.Ctx) error {
	inst := c.Locals("instance").(*instance.Instance)

	quantity := 1
	if q := c.Query("quantity"); q != "" {
		var parsed int
		_, err := fmt.Sscan(q, &parsed)
		if err != nil || parsed == 0 || parsed > 2000 {
			return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Invalid quantity parameter", nil))
		}
		quantity = parsed
	}

	ids, err := h.messageService.GetMessageIDs(context.Background(), inst, input.GenerateMessageIDs{Quantity: quantity})
	if err != nil {
		appErr := app.TranslateError("message handler", err)
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to generate message IDs", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Message IDs generated successfully", fiber.Map{
		"ids": ids,
	}))
}

func (h *MessageHandler) SendText(c fiber.Ctx) error {
	inst := c.Locals("instance").(*instance.Instance)
	var req input.SendTextMessageInput

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	msg, err := h.messageService.SendTextMessage(context.Background(), inst, req)
	if err != nil {
		appErr := app.TranslateError("message handler", err)
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to send message", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Message sent successfully", fiber.Map{
		"message": msg,
	}))
}

func (h *MessageHandler) SendImage(c fiber.Ctx) error {
	inst := c.Locals("instance").(*instance.Instance)
	var req input.SendImageMessageInput

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	msg, err := h.messageService.SendImageMessage(context.Background(), inst, req)
	if err != nil {
		appErr := app.TranslateError("message handler", err)
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to send message", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Message sent successfully", fiber.Map{
		"message": msg,
	}))
}

func (h *MessageHandler) SendVideo(c fiber.Ctx) error {
	inst := c.Locals("instance").(*instance.Instance)
	var req input.SendVideoMessageInput

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	msg, err := h.messageService.SendVideoMessage(context.Background(), inst, req)
	if err != nil {
		appErr := app.TranslateError("message handler", err)
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to send message", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Message sent successfully", fiber.Map{
		"message": msg,
	}))
}

func (h *MessageHandler) SendAudio(c fiber.Ctx) error {
	inst := c.Locals("instance").(*instance.Instance)
	var req input.SendAudioMessageInput

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	msg, err := h.messageService.SendAudioMessage(context.Background(), inst, req)
	if err != nil {
		appErr := app.TranslateError("message handler", err)
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to send message", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Message sent successfully", fiber.Map{
		"message": msg,
	}))
}

func (h *MessageHandler) SendVoice(c fiber.Ctx) error {
	inst := c.Locals("instance").(*instance.Instance)
	var req input.SendVoiceMessageInput

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	msg, err := h.messageService.SendVoiceMessage(context.Background(), inst, req)
	if err != nil {
		appErr := app.TranslateError("message handler", err)
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to send message", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Message sent successfully", fiber.Map{
		"message": msg,
	}))
}

func (h *MessageHandler) SendDocument(c fiber.Ctx) error {
	inst := c.Locals("instance").(*instance.Instance)
	var req input.SendDocumentMessageInput

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	msg, err := h.messageService.SendDocumentMessage(context.Background(), inst, req)
	if err != nil {
		appErr := app.TranslateError("message handler", err)
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to send message", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Message sent successfully", fiber.Map{
		"message": msg,
	}))
}

func (h *MessageHandler) SendReaction(c fiber.Ctx) error {
	inst := c.Locals("instance").(*instance.Instance)
	var req requests.SendReactionMessageInput

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	msg, err := h.messageService.SendReaction(context.Background(), inst, req.ToInput())
	if err != nil {
		appErr := app.TranslateError("message handler", err)
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to send message", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Reaction sent successfully", fiber.Map{
		"reaction": msg,
	}))
}

func (h *MessageHandler) MarkMessagesAsRead(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)

	var req requests.ReadMessagesRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	if bag := req.Validate(); !bag.IsEmpty() {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewValidationErrorResponse(bag))
	}

	err := h.messageService.MarkMessagesAsRead(ctx, inst, req.ToInput())
	if err != nil {
		appErr := app.TranslateError("message handler", err)
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to mark messages as read", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Messages marked as read successfully", nil))
}
