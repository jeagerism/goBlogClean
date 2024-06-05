package miduse

import (
	"errors"

	middlewareRepository "github.com/jeagerism/goBlogClean/modules/middlewares/mid_repo"
)

type middlewareUsecase struct {
	middlewareRepo middlewareRepository.IMiddlewareRepository
}

var ErrUnauthorized = errors.New("unauthorized access")

type IMiddlewareUsecase interface {
	CheckUserRole(userId string) (string, error)
}

func NewMiddlewareUsecase(middlewareRepo middlewareRepository.IMiddlewareRepository) IMiddlewareUsecase {
	return &middlewareUsecase{
		middlewareRepo: middlewareRepo,
	}
}

func (u *middlewareUsecase) CheckUserRole(userId string) (string, error) {
	role, err := u.middlewareRepo.GetUserRole(userId)
	if err != nil {
		return "", err
	}
	if !role {
		return "you are user", ErrUnauthorized
	}
	return "admin", nil
}
