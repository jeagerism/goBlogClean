package usersusecases

import (
	"errors"

	"github.com/jeagerism/goBlogClean/modules/users"
	usersrepositories "github.com/jeagerism/goBlogClean/modules/users/usersRepositories"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

type usersUsecases struct {
	userRepo usersrepositories.IUserRepositories
}

type IUsersUsecases interface {
	Signup(req *users.SignupRequest) (*users.User, error)
	Login(req *users.LoginRequest) (*users.User, error)
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

func (u *usersUsecases) Login(req *users.LoginRequest) (*users.User, error) {
	user, err := u.userRepo.GetUser(req)
	if err != nil {
		return nil, ErrUserNotFound // ให้ข้อมูลข้อผิดพลาดที่ชัดเจน
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidPassword // ให้ข้อมูลข้อผิดพลาดที่ชัดเจน
	}

	return user, nil
}
