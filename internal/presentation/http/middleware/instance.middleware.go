// core/middleware/instance.middleware.go
package middleware

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/service"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http"
)

type InstanceMiddleware struct {
	registry     instance.InstanceRegistry
	instanceRepo instance.InstanceRepository
	instService  *service.SessionService
}

func NewInstanceMiddleware(
	registry instance.InstanceRegistry,
	instanceRepo instance.InstanceRepository,
	instService *service.SessionService,
) *InstanceMiddleware {
	return &InstanceMiddleware{
		registry:     registry,
		instanceRepo: instanceRepo,
		instService:  instService,
	}
}

func (m *InstanceMiddleware) AttachInstance() fiber.Handler {
	l := app.GetMiddlewareLogger()

	return func(c fiber.Ctx) error {
		id := c.Get(http.HeaderInstanceID)
		if id == "" {
			if v, ok := c.Locals("instance_id").(string); ok && v != "" {
				id = v
			} else {
				l.Warn("Missing instance ID")
				appErr := app.NewAppError("instance middleware", app.CodeMissingData, nil)
				return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Missing instance ID in header", appErr))
			}
		}

		inst, ok := m.registry.Get(id)

		if !ok || inst == nil {
			inst, err := m.instanceRepo.Get(instance.WhereID(id))
			if err != nil {
				l.Error("Error on finding instance", "error", err)
				appErr := app.NewAppError("instance middleware", app.CodeInternalError, err)
				return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Error on finding instance", appErr))
			}

			if inst == nil {
				l.Warn("Instance not found", "instance", id)
				appErr := app.NewAppError("instance middleware", app.CodeInstanceNotFound, nil)
				return c.Status(fiber.StatusNotFound).JSON(http.NewErrorResponse("Instance not found", appErr))
			}

			m.registry.Add(inst)
			c.Locals("instance", inst)
			l.Info("Instance loaded and attached to context", "instance", id)
			return c.Next()
		} else {
			l.Debug("Instance found in registry", "instance", id)
		}

		c.Locals("instance", inst)
		return c.Next()
	}
}

func (m *InstanceMiddleware) ConnectInstance() fiber.Handler {
	return func(c fiber.Ctx) error {
		l := app.GetMiddlewareLogger()
		ctx := context.Background()

		val := c.Locals("instance")
		if val == nil {
			l.Warn("Instance not loaded in context")
			appErr := app.NewAppError("instance middleware", app.CodeInstanceNotFound, fmt.Errorf("instance not loaded in context"))
			return c.Status(fiber.StatusInternalServerError).JSON(http.NewErrorResponse("Instance not loaded", appErr))
		}

		inst, ok := val.(*instance.Instance)
		if !ok {
			l.Error("Invalid instance type in context")
			appErr := app.NewAppError("instance middleware", app.CodeInternalError, fmt.Errorf("invalid instance type in context"))
			return c.Status(fiber.StatusInternalServerError).JSON(http.NewErrorResponse("Invalid instance type", appErr))
		}

		if !inst.Status.IsConnected() {
			if err := m.instService.Connect(ctx, inst); err != nil {
				l.Error("Failed to connect instance", "instance", inst.ID, "error", err)
				appErr := app.NewAppError("instance middleware", app.CodeInternalError, err)
				return c.Status(fiber.StatusInternalServerError).JSON(http.NewErrorResponse("Failed to connect instance", appErr))
			}
		}

		l.Info("Instance connected", "instance", inst.ID)
		return c.Next()
	}
}
