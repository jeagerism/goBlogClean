package miduse

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/jeagerism/goBlogClean/internal/jwtclaims"
)

type middlewareUsecase struct {
	jwtSecret []byte
}

type IMiddlewareUsecase interface {
	ParseAccessToken(tokenString string) (*jwtclaims.AccessClaims, error)
}

func NewMiddlewareUsecase(jwtSecret string) IMiddlewareUsecase {
	return &middlewareUsecase{
		jwtSecret: []byte(jwtSecret),
	}
}

func (u *middlewareUsecase) ParseAccessToken(tokenString string) (*jwtclaims.AccessClaims, error) {
	claims := &jwtclaims.AccessClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return u.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}
