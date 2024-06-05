package usershandlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jeagerism/goBlogClean/modules/users"
	usersusecases "github.com/jeagerism/goBlogClean/modules/users/usersUsecases"
)

type usersHandlers struct {
	userUse usersusecases.IUsersUsecases
}

type IUsersHandlers interface {
	Signup(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

func NewUsersHandlers(userUse usersusecases.IUsersUsecases) IUsersHandlers {
	return &usersHandlers{
		userUse: userUse,
	}
}

func (h *usersHandlers) Signup(c *fiber.Ctx) error {
	req := new(users.SignupRequest)
	err := c.BodyParser(req)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user, err := h.userUse.Signup(req)
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{"error": "Signup failed"})
	}
	c.Locals("userId", user.Id)

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *usersHandlers) Login(c *fiber.Ctx) error {
	req := new(users.LoginRequest)
	err := c.BodyParser(req)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user, err := h.userUse.Login(req)
	if err != nil {
		if err == usersusecases.ErrUserNotFound || err == usersusecases.ErrInvalidPassword {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid username or password"})
		}
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{"error": "Login failed"})
	}
	c.Locals("userId", user.Id)

	return c.Status(fiber.StatusOK).JSON(user)
}
