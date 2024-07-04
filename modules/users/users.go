package users

type User struct {
	Id       string `db:"user_id" json:"user_id"`
	UserName string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Role     bool   `db:"role" json:"role"`
}

type SignupRequest struct {
	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Role     bool   `db:"role" json:"role"`
}

type LoginRequest struct {
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}
