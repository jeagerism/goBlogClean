package usersusecases

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	userRepo usersrepositories.IUserRepositories
}

type IUsersUsecases interface {
	Signup(req *users.SignupRequest) (*users.User, error)
	Login(req *users.LoginRequest) (*users.User, string, error)
}

func NewUsersUsecases(userRepo usersrepositories.IUserRepositories) IUsersUsecases {
	return &usersUsecases{
		userRepo: userRepo,
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
		return nil, "", ErrUserNotFound // ให้ข้อมูลข้อผิดพลาดที่ชัดเจน
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, "", ErrInvalidPassword // ให้ข้อมูลข้อผิดพลาดที่ชัดเจน
	}

	// Create the claims
	claims := jwt.MapClaims{
		"username": req.Username,
		"exp":      time.Now().Add(time.Minute * 5).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response
	t, err := token.SignedString([]byte("secret-key"))
	if err != nil {
		return nil, "", ErrGenToken
	}

	return user, t, nil
}
