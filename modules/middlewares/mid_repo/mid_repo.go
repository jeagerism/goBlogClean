package midrepo

import "github.com/jmoiron/sqlx"

type middlewareRepository struct {
	db *sqlx.DB
}

type IMiddlewareRepository interface {
	GetUserRole(userID string) (bool, error)
}

func NewMiddlewareRepository(db *sqlx.DB) IMiddlewareRepository {
	return &middlewareRepository{
		db: db,
	}
}

func (repo *middlewareRepository) GetUserRole(userID string) (bool, error) {
	var role bool
	query := "SELECT role FROM users WHERE user_id = $1"
	err := repo.db.Get(&role, query, userID)
	if err != nil {
		return false, err
	}
	return role, nil
}
