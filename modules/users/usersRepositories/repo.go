package usersrepositories

import (
	"github.com/jeagerism/goBlogClean/modules/users"
	"github.com/jmoiron/sqlx"
)

type userRepositories struct {
	db *sqlx.DB
}

type IUserRepositories interface {
	CreateUser(req *users.SignupRequest) (*users.User, error)
}

func NewUserRepositories(db *sqlx.DB) IUserRepositories {
	return &userRepositories{
		db: db,
	}
}

// Signup POST New
func (r *userRepositories) CreateUser(req *users.SignupRequest) (*users.User, error) {
	query := "INSERT INTO users(username,email,password,role) VALUES ($1,$2,$3,$4) RETURNING user_id"
	var id string
	err := r.db.QueryRow(query, req.Username, req.Email, req.Password, req.Role).Scan(&id)
	if err != nil {
		return nil, err
	}

	user := &users.User{
		Id:       id,
		UserName: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}
	return user, nil
}

//GetUserById

//Login
