package miduse

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	middlewareRepository "github.com/jeagerism/goBlogClean/modules/middlewares/mid_repo"
)

type middlewareUsecase struct {
	middlewareRepo middlewareRepository.IMiddlewareRepository
}

var (
	ErrUnauthorized = errors.New("unauthorized access")
	ErrInvalidToken = errors.New("invalid token")
)

type IMiddlewareUsecase interface {
	CheckUserRole(userId string) (string, error)
	VerifyToken(tokenString string) error
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

func (u *middlewareUsecase) VerifyToken(tokenString string) error {
	//ฟังก์ชัน callback ที่ส่งไปให้ jwt.Parse จะคืนค่า secret key ที่ใช้ในการตรวจสอบลายเซ็นของ token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret-key"), nil
	})

	if err != nil {
		return err
	}
	//ตรวจสอบว่า token นั้นถูกต้องหรือไม่ (token.Valid)
	if !token.Valid {
		return ErrInvalidToken
	}

	return nil
}
