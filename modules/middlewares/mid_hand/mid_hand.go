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

func (h *middlewareHandler) CheckRole() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ตรวจสอบว่า userId ถูกตั้งค่าใน context หรือไม่
		userId := c.Get("userId")
		if userId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized", "userId": userId})
		}

		// ตรวจสอบ role ของผู้ใช้
		role, err := h.middlewareUsecase.CheckUserRole(userId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}

		// ตรวจสอบว่า role ตรงกับ requiredRole หรือไม่
		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden access"})
		}

		// ถ้า role ถูกต้อง ให้ดำเนินการต่อ
		return c.Next()
	}
}

func (h *middlewareHandler) CheckToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or malformed JWT",
			})
		}

		// Check if the token string starts with "Bearer"
		if !strings.HasPrefix(tokenString, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or malformed JWT",
			})
		}

		// Remove "Bearer " prefix
		tokenString = tokenString[len("Bearer "):]

		// Verify the token
		err := h.middlewareUsecase.VerifyToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired JWT"})
		}
		return c.Next()
	}
}
