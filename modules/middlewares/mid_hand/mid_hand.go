package midhand

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	middlewareUsecase "github.com/jeagerism/goBlogClean/modules/middlewares/mid_use"
)

type middlewareHandler struct {
	middlewareUsecase middlewareUsecase.IMiddlewareUsecase
}

type IMiddlewareHandler interface {
	CheckRole() fiber.Handler
	CheckToken() fiber.Handler
}

func NewMiddlewareHandler(middlewareUsecase middlewareUsecase.IMiddlewareUsecase) IMiddlewareHandler {
	return &middlewareHandler{
		middlewareUsecase: middlewareUsecase,
	}
}

// CheckRole enforces admin access using role embedded in the JWT (set by CheckToken).
func (h *middlewareHandler) CheckRole() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, uok := c.Locals("userId").(string)
		if !uok || strings.TrimSpace(uid) == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		role, rok := c.Locals("role").(bool)
		if !rok || !role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden access"})
		}

		return c.Next()
	}
}

// CheckToken validates the Bearer JWT and stores userId and role in Locals for downstream handlers.
func (h *middlewareHandler) CheckToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or malformed JWT",
			})
		}

		if !strings.HasPrefix(tokenString, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or malformed JWT",
			})
		}

		tokenString = tokenString[len("Bearer "):]

		claims, err := h.middlewareUsecase.ParseAccessToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired JWT"})
		}

		c.Locals("userId", claims.UserID)
		c.Locals("role", claims.Role)
		c.Locals("username", claims.Username)

		return c.Next()
	}
}
