package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/service"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/token"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http"
)

type AuthMiddleware struct {
	adminToken   string
	tokenService *service.TokenService
}

func NewAuthMiddleware(adminToken string, tokenService *service.TokenService) *AuthMiddleware {
	return &AuthMiddleware{
		adminToken:   adminToken,
		tokenService: tokenService,
	}
}

func (m *AuthMiddleware) Authenticate() fiber.Handler {
	ctx := context.TODO()
	l := app.GetMiddlewareLogger()

	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			l.Warn("Authorization header is missing")
			appErr := app.NewAppError("auth middleware", app.CodeMissingData, token.ErrInvalidToken)
			return c.Status(fiber.StatusUnauthorized).JSON(http.NewErrorResponse("Missing authorization header", appErr))
		}

		headerToken := strings.TrimPrefix(authHeader, "Bearer ")
		if headerToken == "" {
			l.Warn("Authorization header is invalid")
			appErr := app.NewAppError("auth middleware", app.CodeMissingData, token.ErrInvalidToken)
			return c.Status(fiber.StatusUnauthorized).JSON(http.NewErrorResponse("Invalid authorization header", appErr))
		}

		if headerToken == m.adminToken {
			l.Info("Admin token authenticated")
			c.Locals("is_admin", true)
			return c.Next()
		}

		token, err := m.tokenService.ValidateFullToken(ctx, headerToken)

		if err != nil {
			l.Warn("Token validation error", "error", err)
			appErr := app.NewAppError("auth middleware", app.CodeInvalidToken, err)
			return c.Status(fiber.StatusUnauthorized).JSON(http.NewErrorResponse("Invalid token", appErr))
		}

		l.Info("Token valid", "instance", token.InstanceID)

		c.Locals("is_admin", false)
		c.Locals("instance_id", token.InstanceID)
		return c.Next()
	}
}

func (m *AuthMiddleware) IsAdmin() fiber.Handler {
	l := app.GetMiddlewareLogger()

	return func(c fiber.Ctx) error {
		isAdmin, ok := c.Locals("is_admin").(bool)
		if !ok || !isAdmin {
			l.Warn("Access denied: user is not admin")
			arrErr := app.NewAppError("auth middleware", app.CodeNotAdmin, token.ErrInvalidToken)
			return c.Status(fiber.StatusForbidden).JSON(http.NewErrorResponse("Forbidden", arrErr))
		}
		return c.Next()
	}
}
