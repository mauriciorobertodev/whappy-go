package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) RegisterRoutes(r fiber.Router) {
	r.Get("/health", h.Check)
}

func (h *HealthHandler) Check(c fiber.Ctx) error {
	return c.JSON(http.NewSuccessResponse("Health check successful", fiber.Map{}))
}
