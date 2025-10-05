package handler

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/app/service"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http/middleware"
)

type ContactHandler struct {
	contactService *service.ContactService
}

func NewContactHandler(contactService *service.ContactService) *ContactHandler {
	return &ContactHandler{
		contactService: contactService,
	}
}

func (h *ContactHandler) RegisterRoutes(r fiber.Router, authMiddleware *middleware.AuthMiddleware, instMiddleware *middleware.InstanceMiddleware) {
	adm := r.Group("/contacts", authMiddleware.Authenticate(), instMiddleware.AttachInstance(), instMiddleware.ConnectInstance())

	adm.Get("/", h.GetContacts)
	adm.Get("/:contact", h.GetContact)
	adm.Post("/check", h.Check)
}

func (h *ContactHandler) Check(c fiber.Ctx) error {
	inst := c.Locals("instance").(*instance.Instance)
	var req struct {
		Phones []string `json:"phones"`
	}

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	checked, err := h.contactService.Check(context.Background(), inst, input.CheckPhones{
		Phones: req.Phones,
	})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to create instance", err))
	}

	return c.JSON(http.NewSuccessResponse("Phones checked successfully", fiber.Map{
		"checked": checked,
	}))
}

func (h *ContactHandler) GetContact(c fiber.Ctx) error {
	inst := c.Locals("instance").(*instance.Instance)
	phoneOrJID := c.Params("contact")

	contact, err := h.contactService.GetContact(context.Background(), inst, input.GetContact{
		PhoneOrJID: phoneOrJID,
	})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to get contact", err))
	}

	return c.JSON(http.NewSuccessResponse("Contact retrieved successfully", fiber.Map{
		"contact": contact,
	}))
}

func (h *ContactHandler) GetContacts(c fiber.Ctx) error {
	inst := c.Locals("instance").(*instance.Instance)

	contacts, err := h.contactService.GetContacts(context.Background(), inst)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to get contacts", err))
	}

	return c.JSON(http.NewSuccessResponse("Contacts retrieved successfully", fiber.Map{
		"contacts": contacts,
	}))
}
