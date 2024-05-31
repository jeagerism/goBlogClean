package blogs

type Blog struct {
	Id        int    `db:"blog_id" json:"blog_id"`
	UserId    int    `db:"user_id" json:"user_id"`
	Title     string `db:"title" json:"title"`
	Content   string `db:"content" json:"content"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

type BlogRequest struct {
	Title   string `db:"title" json:"title"`
	Content string `db:"content" json:"content"`
}
