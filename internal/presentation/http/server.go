package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/mauriciorobertodev/whappy-go/internal/app"
)

const HeaderInstanceID = "X-Instance-ID"

func NewFiberApp(appName, appVersion string, isProduction bool) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: appName + " v" + appVersion,
		ErrorHandler: func(c fiber.Ctx, err error) error {
			if appErr, ok := err.(*app.AppError); ok {
				return c.Status(fiber.StatusInternalServerError).JSON(NewErrorResponse("Internal server error", app.WrapLoc("api error handler", appErr)))
			}

			// TODO: interpret fiber errors

			// Padronização...
			return c.Status(fiber.StatusInternalServerError).JSON(NewInternalErrorResponse("api error handler", err))
		},
		BodyLimit: 150 * 1024 * 1024, // 150 MB
	})

	app.Use(logger.New(logger.Config{
		// TimeZone: "UTC",
	}))

	app.Use(cors.New())

	if isProduction {
		// app.Use(recoverer.New())
	}

	return app
}
