package blogs

import "time"

type Blog struct {
	Id        string    `db:"blog_id" json:"blog_id"`
	UserId    string    `db:"user_id" json:"user_id"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type BlogRequest struct {
	UserId  string `db:"user_id" json:"user_id"`
	Title   string `db:"title" json:"title"`
	Content string `db:"content" json:"content"`
}
