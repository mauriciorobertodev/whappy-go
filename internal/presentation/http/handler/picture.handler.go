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

type PictureHandler struct {
	pictureService *service.PictureService
}

func NewPictureHandler(pictureService *service.PictureService) *PictureHandler {
	return &PictureHandler{
		pictureService: pictureService,
	}
}

func (h *PictureHandler) RegisterRoutes(r fiber.Router, authMiddleware *middleware.AuthMiddleware, instMiddleware *middleware.InstanceMiddleware) {
	adm := r.Group("/pictures", authMiddleware.Authenticate(), instMiddleware.AttachInstance(), instMiddleware.ConnectInstance())

	adm.Get("/:picture", h.Get)
}

func (h *PictureHandler) Get(c fiber.Ctx) error {
	inst := c.Locals("instance").(*instance.Instance)
	picture := c.Params("picture")
	preview := c.Query("preview", "true") == "true"
	isCommunity := c.Query("is_community", "false") == "true"

	inp := input.GetPictureInput{
		PhoneOrJID:  picture,
		Preview:     &preview,
		IsCommunity: &isCommunity,
	}

	pictureURL, err := h.pictureService.Get(context.Background(), inst, inp)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to get picture", err))
	}

	return c.JSON(http.NewSuccessResponse("Picture retrieved successfully", fiber.Map{
		"picture": pictureURL,
	}))
}
