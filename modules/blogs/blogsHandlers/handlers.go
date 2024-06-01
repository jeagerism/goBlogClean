package blogshandlers

import (
	"github.com/gofiber/fiber/v2"
	blogsusecases "github.com/jeagerism/goBlogClean/modules/blogs/blogsUsecases"
)

type blogsHandlers struct {
	blogsUse blogsusecases.IBlogsUsecases
}

type IBlogsHandlers interface {
	FindBlogs(c *fiber.Ctx) error
}

func NewBlogsHandlers(blogsUse blogsusecases.IBlogsUsecases) IBlogsHandlers {
	return &blogsHandlers{
		blogsUse: blogsUse,
	}
}

func (h *blogsHandlers) FindBlogs(c *fiber.Ctx) error {
	blogs, err := h.blogsUse.GetAllBlogs()
	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, "data not found")
	}
	return c.Status(fiber.StatusOK).JSON(blogs)
}
