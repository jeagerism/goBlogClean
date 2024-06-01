package blogsusecases

import (
	"github.com/jeagerism/goBlogClean/modules/blogs"
	blogsrepositories "github.com/jeagerism/goBlogClean/modules/blogs/blogsRepositories"
)

type blogsUseCases struct {
	blogsRepo blogsrepositories.IBlogsRepositories
}

type IBlogsUsecases interface {
	GetAllBlogs() ([]blogs.Blog, error)
}

func NewBlogsUsecase(blogsRepo blogsrepositories.IBlogsRepositories) IBlogsUsecases {
	return &blogsUseCases{
		blogsRepo: blogsRepo,
	}
}

func (u *blogsUseCases) GetAllBlogs() ([]blogs.Blog, error) {
	return u.blogsRepo.GetAll()
}
