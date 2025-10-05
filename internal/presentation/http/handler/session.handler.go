package handler

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/service"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http/middleware"
)

type SessionHandler struct {
	sessionService *service.SessionService
}

func NewSessionHandler(sessionService *service.SessionService) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
	}
}

func (h *SessionHandler) RegisterRoutes(r fiber.Router, authMiddleware *middleware.AuthMiddleware, instMiddleware *middleware.InstanceMiddleware) {
	sessionRouter := r.Group("/session", authMiddleware.Authenticate(), instMiddleware.AttachInstance())
	sessionRouter.Get("/ping", h.Ping)
	sessionRouter.Post("/login", h.Pair)
	sessionRouter.Delete("/logout", h.Logout)
	sessionRouter.Post("/connect", h.Connect)
	sessionRouter.Delete("/disconnect", h.Disconnect)
	sessionRouter.Get("/qr", h.QrCode)
}

func (h *SessionHandler) Ping(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)

	state, err := h.sessionService.Ping(ctx, inst)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(http.NewInternalErrorResponse("session.handler", err))
	}

	return c.JSON(http.NewSuccessResponse("Instance pinged successfully", fiber.Map{
		"state": state,
	}))
}

func (h *SessionHandler) Pair(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)

	if err := h.sessionService.Pair(ctx, inst); err != nil {
		switch err.Code {
		case app.CodeInstanceAlreadyLoggedIn:
			return c.Status(fiber.StatusConflict).JSON(http.NewErrorResponse("Instance is already logged in", err))
		case app.CodeInstanceIsPairing:
			return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Instance is in pairing mode", err))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(http.NewInternalErrorResponse("session.handler", err))
	}

	return c.JSON(http.NewSuccessResponse("Instance logged in successfully", nil))
}

func (h *SessionHandler) Logout(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)

	if err := h.sessionService.Logout(ctx, inst); err != nil {
		switch err.Code {
		case app.CodeInstanceAlreadyLoggedOut:
			return c.Status(fiber.StatusConflict).JSON(http.NewErrorResponse("Instance is already logged out", err))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(http.NewInternalErrorResponse("session.handler", err))
	}

	return c.JSON(http.NewSuccessResponse("Instance logged out successfully", nil))
}

func (h *SessionHandler) Connect(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)

	if err := h.sessionService.Connect(ctx, inst); err != nil {
		switch err.Code {
		case app.CodeInstanceNotLoggedIn:
			return c.Status(fiber.StatusUnauthorized).JSON(http.NewErrorResponse("Instance is not logged in", err))
		case app.CodeInstanceBanned:
			return c.Status(fiber.StatusLocked).JSON(http.NewErrorResponse("Instance is permanently banned", err))
		case app.CodeInstanceAlreadyConnected:
			return c.Status(fiber.StatusAlreadyReported).JSON(http.NewErrorResponse("Instance is already connected", err))
		case app.CodeInstanceIsConnecting:
			return c.Status(fiber.StatusAccepted).JSON(http.NewErrorResponse("Instance is in connecting mode", err))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(http.NewInternalErrorResponse("session.handler", err))
	}

	return c.JSON(http.NewSuccessResponse("Instance connected successfully", nil))
}

func (h *SessionHandler) Disconnect(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)

	if err := h.sessionService.Disconnect(ctx, inst); err != nil {
		switch err.Code {
		case app.CodeInstanceAlreadyDisconnected:
			return c.Status(fiber.StatusConflict).JSON(http.NewErrorResponse("Instance is already disconnected", err))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(http.NewInternalErrorResponse("session.handler", err))
	}

	return c.JSON(http.NewSuccessResponse("Instance disconnected successfully", nil))
}

func (h *SessionHandler) QrCode(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)

	qr, err := h.sessionService.QrCode(ctx, inst)
	if err != nil {
		switch err.Code {
		case app.CodeInstanceAlreadyPaired:
			return c.Status(fiber.StatusConflict).JSON(http.NewErrorResponse("Instance is already paired", err))
		case app.CodeInstanceNotPairing:
			return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Instance is not in pairing mode", err))
		case app.CodeInstanceNoQrCode:
			return c.Status(fiber.StatusNotFound).JSON(http.NewErrorResponse("No QR code found for instance, wait some time and try again", err))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(http.NewInternalErrorResponse("session.handler", err))
	}

	return c.JSON(http.NewSuccessResponse("QR code retrieved successfully", fiber.Map{
		"qr_code": qr,
	}))
}
