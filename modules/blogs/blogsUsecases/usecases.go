package blogsusecases

import (
	"github.com/jeagerism/goBlogClean/modules/blogs"
	blogsrepositories "github.com/jeagerism/goBlogClean/modules/blogs/blogsRepositories"
)

type blogsUseCases struct {
	blogsRepo blogsrepositories.IBlogsRepositories
}

type IBlogsUsecases interface {
	GetAllBlogs(page, limit int) ([]blogs.Blog, *blogs.Pagination, error)
	GetBlogById(id string) (*blogs.Blog, error)
	PostBlog(req *blogs.BlogRequest) (*blogs.Blog, error)
	UpdateBlog(req *blogs.BlogUpdateRequest) (*blogs.Blog, error)
	DeleteBlog(id string) error
}

func NewBlogsUsecase(blogsRepo blogsrepositories.IBlogsRepositories) IBlogsUsecases {
	return &blogsUseCases{
		blogsRepo: blogsRepo,
	}
}

func (u *blogsUseCases) GetAllBlogs(page, limit int) ([]blogs.Blog, *blogs.Pagination, error) {
	blogs, paginate, err := u.blogsRepo.GetAll(page, limit)
	if err != nil {
		return nil, nil, err
	}
	return blogs, paginate, nil
}

func (u *blogsUseCases) GetBlogById(id string) (*blogs.Blog, error) {
	blog, err := u.blogsRepo.GetById(id)
	if err != nil {
		return nil, err
	}
	return blog, nil
}

func (u *blogsUseCases) PostBlog(req *blogs.BlogRequest) (*blogs.Blog, error) {
	blog, err := u.blogsRepo.Post(req)
	if err != nil {
		return nil, err
	}
	return blog, nil
}

func (u *blogsUseCases) UpdateBlog(req *blogs.BlogUpdateRequest) (*blogs.Blog, error) {
	blog, err := u.blogsRepo.Update(req)
	if err != nil {
		return nil, err
	}
	return blog, nil
}

func (u *blogsUseCases) DeleteBlog(id string) error {
	err := u.blogsRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
