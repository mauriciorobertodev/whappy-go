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
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http/requests"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http/resources"
)

type InstanceHandler struct {
	instService  *service.InstanceService
	instRegistry instance.InstanceRegistry
}

func NewInstanceHandler(instService *service.InstanceService, instRegistry instance.InstanceRegistry) *InstanceHandler {
	return &InstanceHandler{
		instService:  instService,
		instRegistry: instRegistry,
	}
}

func (h *InstanceHandler) RegisterRoutes(r fiber.Router, auth *middleware.AuthMiddleware) {
	adm := r.Group("/instances", auth.Authenticate(), auth.IsAdmin())

	adm.Post("", h.Create)
	adm.Get("", h.List)
	adm.Get("/:id", h.Get)
	adm.Put("/:id/token", h.Token)
}

func (h *InstanceHandler) List(c fiber.Ctx) error {
	instances, err := h.instService.List(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(http.NewErrorResponse("Failed to retrieve instances", err))
	}

	return c.Status(fiber.StatusOK).JSON(http.NewSuccessResponse("Instances retrieved successfully", fiber.Map{
		"instances": resources.MakeInstanceResources(instances),
	}))
}

func (h *InstanceHandler) Create(c fiber.Ctx) error {
	var req requests.CreateInstanceRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewInvalidJSONResponse())
	}

	if bag := req.Validate(); bag.HasErrors() {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(http.NewValidationErrorResponse(bag))
	}

	inst, token, err := h.instService.Create(context.Background(), req.ToInput())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to create instance", err))
	}

	return c.JSON(http.NewSuccessResponse("Instance created successfully", fiber.Map{
		"instance": resources.MakeInstanceResource(inst),
		"token":    token.FullToken(),
	}))
}

func (h *InstanceHandler) Get(c fiber.Ctx) error {
	id := c.Params("id")

	instance, err := h.instService.Get(context.Background(), input.GetInstance{ID: id})
	if err != nil {
		if err.Code == app.CodeInstanceNotFound || instance == nil {
			return c.Status(fiber.StatusNotFound).JSON(http.NewErrorResponse("Instance not found", err))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(http.NewErrorResponse("Failed to retrieve instance", err))
	}

	return c.JSON(http.NewSuccessResponse("Instance retrieved successfully", fiber.Map{
		"instance": resources.MakeInstanceResource(instance),
	}))
}

func (h *InstanceHandler) Token(c fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	inst, err := h.instService.Get(ctx, input.GetInstance{ID: id})
	if err != nil {
		if err.Code == app.CodeInstanceNotFound || inst == nil {
			return c.Status(fiber.StatusNotFound).JSON(http.NewErrorResponse("Instance not found", err))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(http.NewErrorResponse("Failed to retrieve instance", err))
	}

	tok, appErr := h.instService.RenewToken(ctx, input.RenewInstanceToken{ID: id})
	if appErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(http.NewErrorResponse("Failed to renew token", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Token renewed successfully", fiber.Map{
		"token": tok.FullToken(),
	}))
}
