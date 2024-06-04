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
		return c.Status(fiber.ErrBadRequest.Code).JSON("body error")
	}

	user, err := h.userUse.Signup(req)
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON("signup failed")
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}
