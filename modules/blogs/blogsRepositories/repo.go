package blogsrepositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/jeagerism/goBlogClean/modules/blogs"
	"github.com/jmoiron/sqlx"
)

// ========>>>> 3 Comboset
type IBlogsRepositories interface {
	GetAll(page, limit int) ([]blogs.Blog, *blogs.Pagination, error)
	GetById(id string) (*blogs.Blog, error)
	Post(req *blogs.BlogRequest) (*blogs.Blog, error)
	Update(req *blogs.BlogUpdateRequest) (*blogs.Blog, error)
	Delete(id string) error
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
func (r *blogsRepositories) GetAll(page, limit int) ([]blogs.Blog, *blogs.Pagination, error) {
	var blogs []blogs.Blog
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := limit * (page - 1)
	query := "SELECT * FROM blogs ORDER BY blog_id LIMIT $1 OFFSET $2"
	err := r.db.Select(&blogs, query, limit, offset)
	if err != nil {
		return nil, nil, err
	}

	pagination := pagination(r.db, "blogs", limit, page)
	return blogs, pagination, nil
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

// Update
func (r *blogsRepositories) Update(req *blogs.BlogUpdateRequest) (*blogs.Blog, error) {
	query := "UPDATE blogs SET title = $1, content = $2 WHERE blog_id = $3"
	result, err := r.db.Exec(query, req.Title, req.Content, req.Id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, errors.New("no rows affected")
	}
	blog, err := r.GetById(req.Id)
	if err != nil {
		return nil, err
	}

	return blog, nil
}

// Delete
func (r *blogsRepositories) Delete(id string) error {
	query := "DELETE FROM blogs WHERE blog_id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffeted, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffeted == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

// Pagination
func pagination(db *sqlx.DB, table string, limit, page int) *blogs.Pagination {
	var tmpl blogs.Pagination
	var recordCount int

	// คำสั่ง SQL สำหรับการนับจำนวนเรคคอร์ดทั้งหมดในตารางที่ระบุ
	sqlCount := fmt.Sprintf("SELECT count(blog_id) FROM %s", table)

	// ดึงจำนวนเรคคอร์ดทั้งหมดจากฐานข้อมูล
	db.QueryRow(sqlCount).Scan(&recordCount)

	// คำนวณหน้าทั้งหมด
	total := (recordCount / limit)

	reminder := (recordCount % limit)
	if reminder == 0 {
		tmpl.TotalPage = total
	} else {
		tmpl.TotalPage = total + 1
	}

	// ตั้งค่าหน้าปัจจุบันและจำนวนรายการต่อหน้า
	tmpl.CurrentPage = page
	tmpl.RecordPerPage = limit

	// คำนวณหน้าถัดไปและหน้าก่อนหน้า
	if page <= 0 {
		tmpl.Next = 1
	} else if page < tmpl.TotalPage {
		tmpl.Previous = page - 1
		tmpl.Next = page + 1
	} else if page == tmpl.TotalPage {
		tmpl.Previous = page - 1
		tmpl.Next = 0
	}

	return &tmpl
}
