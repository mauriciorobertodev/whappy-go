package handler

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/app/service"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http/middleware"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http/requests"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http/resources"
)

type WebhookHandler struct {
	webhookService *service.WebhookService
}

func NewWebhookHandler(webhookService *service.WebhookService) *WebhookHandler {
	return &WebhookHandler{
		webhookService: webhookService,
	}
}

func (h *WebhookHandler) RegisterRoutes(r fiber.Router, authMiddleware *middleware.AuthMiddleware, instMiddleware *middleware.InstanceMiddleware) {
	we := r.Group("/webhooks", authMiddleware.Authenticate(), instMiddleware.AttachInstance())

	we.Get("/", h.GetWebhooks)
	we.Post("/", h.CreateWebhook)
	we.Get("/:id", h.GetWebhook)
	we.Put("/:id", h.UpdateWebhook)
	we.Delete("/:id", h.DeleteWebhook)
}

func (h *WebhookHandler) GetWebhook(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)
	id := c.Params("id")

	webhook, appErr := h.webhookService.GetWebhook(ctx, inst, input.GetWebhook{
		ID: id,
	})

	if appErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to get webhook", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Webhook retrieved successfully", fiber.Map{
		"webhook": resources.MakeWebhookResource(webhook, nil),
	}))
}

func (h *WebhookHandler) GetWebhooks(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)

	webhooks, appErr := h.webhookService.GetWebhooks(ctx, inst)
	if appErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Invalid limit parameter", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Webhooks retrieved successfully", fiber.Map{
		"webhooks": webhooks,
	}))
}

func (h *WebhookHandler) CreateWebhook(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)
	var req requests.CreateWebhook
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	webhook, secret, appErr := h.webhookService.CreateWebhook(ctx, inst, req.ToInput())
	if appErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to create webhook", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Webhook created successfully", fiber.Map{
		"webhook": resources.MakeWebhookResource(webhook, &secret),
	}))
}

func (h *WebhookHandler) UpdateWebhook(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)
	id := c.Params("id")
	var req requests.UpdateWebhook
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	webhook, appErr := h.webhookService.UpdateWebhook(ctx, inst, req.ToInput(id))
	if appErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to update webhook", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Webhook updated successfully", fiber.Map{
		"webhook": webhook,
	}))
}

func (h *WebhookHandler) DeleteWebhook(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)
	id := c.Params("id")

	appErr := h.webhookService.DeleteWebhook(ctx, inst, input.DeleteWebhook{
		ID: id,
	})

	if appErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to delete webhook", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Webhook deleted successfully", nil))
}
