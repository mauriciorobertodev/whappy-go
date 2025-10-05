package handler

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/app/service"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http/middleware"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

func (h *ChatHandler) RegisterRoutes(r fiber.Router, authMiddleware *middleware.AuthMiddleware, instMiddleware *middleware.InstanceMiddleware) {
	msg := r.Group("/chat", authMiddleware.Authenticate(), instMiddleware.AttachInstance(), instMiddleware.ConnectInstance())
	msg.Post("/presence", h.Presence)
}

func (h *ChatHandler) Presence(c fiber.Ctx) error {
	inst := c.Locals("instance").(*instance.Instance)
	var req input.SendChatPresenceInput

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	if err := h.chatService.SendPresence(context.Background(), inst, req); err != nil {
		appErr := app.TranslateError("chat handler", err)
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to send presence", appErr))
	}

	return c.Status(fiber.StatusNoContent).JSON(http.NewSuccessResponse("Presence sent successfully", nil))
}
