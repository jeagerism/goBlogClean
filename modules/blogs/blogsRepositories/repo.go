package blogsrepositories

import (
	"time"

	"github.com/jeagerism/goBlogClean/modules/blogs"
	"github.com/jmoiron/sqlx"
)

// ========>>>> 3 Comboset
type IBlogsRepositories interface {
	GetAll() ([]blogs.Blog, error)
	GetById(id string) (*blogs.Blog, error)
	Post(req *blogs.BlogRequest) (*blogs.Blog, error)
}

type blogsRepositories struct {
	db *sqlx.DB
}

func NewBlogsRepositories(db *sqlx.DB) IBlogsRepositories {
	return &blogsRepositories{
		db: db,
	}
}

// ========>>>> 3 Comboset

// *[] for check nil [] for check len = 0
func (r *blogsRepositories) GetAll() ([]blogs.Blog, error) {
	var blogs []blogs.Blog
	query := "SELECT * FROM blogs"
	err := r.db.Select(&blogs, query)
	if err != nil {
		return nil, err
	}
	return blogs, nil
}

func (r *blogsRepositories) GetById(id string) (*blogs.Blog, error) {
	var blog blogs.Blog
	query := "SELECT * FROM blogs WHERE blog_id = $1"
	err := r.db.Get(&blog, query, id)
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *blogsRepositories) Post(req *blogs.BlogRequest) (*blogs.Blog, error) {

	query := "INSERT INTO blogs(user_id,title,content) VALUES ($1,$2,$3) RETURNING blog_id,created_at"
	var id string
	var created_at time.Time
	err := r.db.QueryRow(query, req.UserId, req.Title, req.Content).Scan(&id, &created_at)
	if err != nil {
		return nil, err
	}
	blog := blogs.Blog{
		Id:        id,
		UserId:    req.UserId,
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: created_at,
	}
	return &blog, nil
}

//Update

//Delete
