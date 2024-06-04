package usersusecases

import (
	"github.com/jeagerism/goBlogClean/modules/users"
	usersrepositories "github.com/jeagerism/goBlogClean/modules/users/usersRepositories"
	"golang.org/x/crypto/bcrypt"
)

type usersUsecases struct {
	userRepo usersrepositories.IUserRepositories
}

type IUsersUsecases interface {
	Signup(req *users.SignupRequest) (*users.User, error)
}

func NewUsersUsecases(userRepo usersrepositories.IUserRepositories) IUsersUsecases {
	return &usersUsecases{
		userRepo: userRepo,
	}
}

func (u *usersUsecases) Signup(req *users.SignupRequest) (*users.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, err
	}

	req.Password = string(hashedPassword)
	user, err := u.userRepo.CreateUser(req)
	if err != nil {
		return nil, err
	}

	return user, nil
}
