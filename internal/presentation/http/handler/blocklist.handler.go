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

type BlocklistHandler struct {
	blocklistService *service.BlocklistService
}

func NewBlocklistHandler(blocklistService *service.BlocklistService) *BlocklistHandler {
	return &BlocklistHandler{
		blocklistService: blocklistService,
	}
}

func (h *BlocklistHandler) RegisterRoutes(r fiber.Router, authMiddleware *middleware.AuthMiddleware, instMiddleware *middleware.InstanceMiddleware) {
	adm := r.Group("/blocklist", authMiddleware.Authenticate(), instMiddleware.AttachInstance(), instMiddleware.ConnectInstance())

	adm.Get("/", h.GetBlocklist)
	adm.Post("/:contact", h.Block)
	adm.Delete("/:contact", h.Unblock)
}

func (h *BlocklistHandler) GetBlocklist(c fiber.Ctx) error {
	inst := c.Locals("instance").(*instance.Instance)

	blocklist, err := h.blocklistService.GetBlocklist(context.Background(), inst)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to get blocklist", err))
	}

	return c.JSON(http.NewSuccessResponse("Blocklist retrieved successfully", fiber.Map{
		"blocklist": blocklist,
	}))
}

func (h *BlocklistHandler) Block(c fiber.Ctx) error {
	inst := c.Locals("instance").(*instance.Instance)
	phoneOrJID := c.Params("contact")

	blocklist, err := h.blocklistService.Block(context.Background(), inst, input.Block{
		PhoneOrJID: phoneOrJID,
	})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to block contact", err))
	}

	return c.JSON(http.NewSuccessResponse("Contact blocked successfully", fiber.Map{
		"blocklist": blocklist,
	}))
}

func (h *BlocklistHandler) Unblock(c fiber.Ctx) error {
	inst := c.Locals("instance").(*instance.Instance)
	phoneOrJID := c.Params("contact")

	blocklist, err := h.blocklistService.Unblock(context.Background(), inst, input.Unblock{
		PhoneOrJID: phoneOrJID,
	})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to unblock contact", err))
	}

	return c.JSON(http.NewSuccessResponse("Contact unblocked successfully", fiber.Map{
		"blocklist": blocklist,
	}))
}
