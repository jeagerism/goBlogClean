package usersusecases

import (
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jeagerism/goBlogClean/internal/jwtclaims"
	"github.com/jeagerism/goBlogClean/modules/users"
	usersrepositories "github.com/jeagerism/goBlogClean/modules/users/usersRepositories"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrGenToken        = errors.New("could not generate token")
)

type usersUsecases struct {
	userRepo  usersrepositories.IUserRepositories
	jwtSecret string
}

type IUsersUsecases interface {
	Signup(req *users.SignupRequest) (*users.User, error)
	Login(req *users.LoginRequest) (*users.User, string, error)
}

func NewUsersUsecases(userRepo usersrepositories.IUserRepositories, jwtSecret string) IUsersUsecases {
	return &usersUsecases{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (u *usersUsecases) Signup(req *users.SignupRequest) (*users.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	req.Password = string(hashedPassword)
	user, err := u.userRepo.CreateUser(req)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *usersUsecases) Login(req *users.LoginRequest) (*users.User, string, error) {
	user, err := u.userRepo.GetUser(req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", ErrUserNotFound
		}
		return nil, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, "", ErrInvalidPassword
	}

	now := time.Now()
	claims := &jwtclaims.AccessClaims{
		UserID:   user.Id,
		Username: user.UserName,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(60 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   user.Id,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return nil, "", ErrGenToken
	}

	return user, t, nil
}
