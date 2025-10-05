package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/app/service"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/file"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http/middleware"
	"github.com/mauriciorobertodev/whappy-go/internal/utils"
)

type UploadHandler struct {
	uploadService *service.UploadService
}

func NewUploadHandler(uploadService *service.UploadService) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
	}
}

func (h *UploadHandler) RegisterRoutes(r fiber.Router, authMiddleware *middleware.AuthMiddleware, instMiddleware *middleware.InstanceMiddleware) {
	up := r.Group("/uploads", authMiddleware.Authenticate(), instMiddleware.AttachInstance())

	up.Get("/", h.ListUploads)
	up.Post("/", h.UploadFile)
	up.Get("/:id", h.Get)
	up.Delete("/:id", h.Delete)
}

func (h *UploadHandler) ListUploads(c fiber.Ctx) error {
	inst := c.Locals("instance").(*instance.Instance)

	lim := c.Query("limit", "20")
	var limit int
	_, err := fmt.Sscanf(lim, "%d", &limit)

	if err != nil {
		appErr := app.NewAppError("upload handler", app.CodeInvalidJSON, err)
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Invalid limit parameter", appErr))
	}

	cursor := c.Query("cursor", "")

	files, nextCursorEncoded, appErr := h.uploadService.ListUploads(context.Background(), inst, input.ListUploads{
		Cursor: utils.StringPtr(cursor),
		Limit:  limit,
	})

	if appErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(http.NewInternalErrorResponse("upload handler", appErr))
	}

	return c.JSON(http.NewSuccessResponse("Files retrieved successfully", fiber.Map{
		"files":  files,
		"cursor": nextCursorEncoded,
	}))
}

func (h *UploadHandler) UploadFile(c fiber.Ctx) error {
	ctx := context.Background()
	inst := c.Locals("instance").(*instance.Instance)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		appError := app.NewAppError("upload handler", app.CodeMissingData, err)
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("Failed to get uploaded file", appError))
	}

	stream, err := fileHeader.Open()
	if err != nil {
		appError := app.NewAppError("upload handler", app.CodeCorruptedFile, err)
		return c.Status(fiber.StatusInternalServerError).JSON(http.NewErrorResponse("Failed to open uploaded file", appError))
	}

	defer stream.Close()

	name := c.FormValue("name")
	mimeType := c.FormValue("mime")
	width := c.FormValue("width")
	height := c.FormValue("height")
	duration := c.FormValue("duration")
	pages := c.FormValue("pages")
	thumbnailID := c.FormValue("thumbnail_id")

	metadata := file.Metadata{}

	if name != "" {
		metadata.Name = &name
	}
	if mimeType != "" {
		metadata.Mime = &mimeType
	}
	if width != "" {
		metadata.Width = new(uint32)
		fmt.Sscanf(width, "%d", metadata.Width)
	}
	if height != "" {
		metadata.Height = new(uint32)
		fmt.Sscanf(height, "%d", metadata.Height)
	}
	if duration != "" {
		metadata.Duration = new(uint32)
		fmt.Sscanf(duration, "%d", metadata.Duration)
	}
	if pages != "" {
		metadata.Pages = new(uint32)
		fmt.Sscanf(pages, "%d", metadata.Pages)
	}

	file, appErr := h.uploadService.UploadWithStream(ctx, inst, input.UploadFile{
		Stream:      stream,
		Metadata:    metadata,
		ThumbnailID: utils.StringPtr(thumbnailID),
	})
	if appErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(http.NewInternalErrorResponse("upload handler", appErr))
	}

	return c.JSON(http.NewSuccessResponse("File uploaded successfully", fiber.Map{
		"file": file,
	}))
}

func (h *UploadHandler) Get(c fiber.Ctx) error {
	inst := c.Locals("instance").(*instance.Instance)

	id := c.Params("id")
	if id == "" {
		appErr := app.NewAppError("upload handler", app.CodeMissingData, errors.New("file ID is required"))
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("File ID is required", appErr))
	}

	file, appErr := h.uploadService.GetUpload(context.Background(), inst, input.GetUpload{
		FileID: id,
	})

	if appErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(http.NewInternalErrorResponse("upload handler", appErr))
	}

	return c.JSON(http.NewSuccessResponse("File retrieved successfully", fiber.Map{
		"file": file,
	}))
}

func (h *UploadHandler) Delete(c fiber.Ctx) error {
	inst := c.Locals("instance").(*instance.Instance)
	id := c.Params("id")
	if id == "" {
		appErr := app.NewAppError("upload handler", app.CodeMissingData, errors.New("file ID is required"))
		return c.Status(fiber.StatusBadRequest).JSON(http.NewErrorResponse("File ID is required", appErr))
	}

	appErr := h.uploadService.DeleteUpload(context.Background(), inst, input.DeleteUpload{
		FileID: id,
	})

	if appErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(http.NewInternalErrorResponse("upload handler", appErr))
	}

	return c.JSON(http.NewSuccessResponse("File deleted successfully", nil))
}
