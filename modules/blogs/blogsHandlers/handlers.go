package blogshandlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jeagerism/goBlogClean/modules/blogs"
	blogsusecases "github.com/jeagerism/goBlogClean/modules/blogs/blogsUsecases"
)

type blogsHandlers struct {
	blogsUsecases blogsusecases.IBlogsUsecases
}

type IBlogsHandlers interface {
	FindBlogs(c *fiber.Ctx) error
	FindBlog(c *fiber.Ctx) error
	PostBlog(c *fiber.Ctx) error
	UpdateBlog(c *fiber.Ctx) error
	DeleteBlog(c *fiber.Ctx) error
}

func NewBlogsHandlers(blogsUsecases blogsusecases.IBlogsUsecases) IBlogsHandlers {
	return &blogsHandlers{
		blogsUsecases: blogsUsecases,
	}
}

func (h *blogsHandlers) FindBlogs(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)    // ใช้ QueryInt เพื่อแปลงค่าที่ส่งมาเป็น int
	limit := c.QueryInt("limit", 10) // ค่า default คือ 10

	blogs, pagination, err := h.blogsUsecases.GetAllBlogs(page, limit)
	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, "Data not found")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":     blogs,
		"paginate": pagination,
	})
}

func (h *blogsHandlers) FindBlog(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("blogId"))
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid Blog ID",
		})
	}

	blog, err := h.blogsUsecases.GetBlogById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Blog not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(blog)
}

func (h *blogsHandlers) PostBlog(c *fiber.Ctx) error {
	request := new(blogs.BlogRequest)
	if err := c.BodyParser(request); err != nil {
		return fiber.NewError(fiber.ErrBadRequest.Code, "Failed to parse blog content")
	}

	if request.UserId == "" || request.Title == "" || request.Content == "" {
		return c.Status(fiber.ErrUnprocessableEntity.Code).JSON(fiber.Map{
			"error":   true,
			"message": "Missing required blog fields",
		})
	}

	blog, err := h.blogsUsecases.PostBlog(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create blog",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(blog)
}

func (h *blogsHandlers) UpdateBlog(c *fiber.Ctx) error {
	request := new(blogs.BlogUpdateRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to parse blog update request",
		})
	}

	if request.Title == "" || request.Content == "" || request.Id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Missing required blog update fields",
		})
	}

	blog, err := h.blogsUsecases.UpdateBlog(request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to update blog",
		})
	}
	return c.Status(fiber.StatusOK).JSON(blog)
}

func (h *blogsHandlers) DeleteBlog(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("blogId"))
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid Blog ID",
		})
	}

	err := h.blogsUsecases.DeleteBlog(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete blog",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Delete success!",
	})
}
