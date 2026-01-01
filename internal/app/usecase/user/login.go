package user

import (
	"errors"

	"github.com/wellingtonlope/ticket-api/internal/app/repository"
)

var ErrUserEmailPasswordWrong = errors.New("email or password is wrong")

type Login interface {
	Handle(input LoginInput) (*Output, error)
}

type login struct {
	userRepository repository.UserRepository
}

func NewLogin(userRepository repository.UserRepository) Login {
	return &login{userRepository: userRepository}
}

type LoginInput struct {
	Email    string
	Password string
}

func (u *login) Handle(input LoginInput) (*Output, error) {
	user, err := u.userRepository.GetByEmail(input.Email)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, ErrUserEmailPasswordWrong
		}
		return nil, err
	}

	if !user.Password.IsCorrectPassword(input.Password) {
		return nil, ErrUserEmailPasswordWrong
	}

	return userOutputFromUser(user), nil
}
