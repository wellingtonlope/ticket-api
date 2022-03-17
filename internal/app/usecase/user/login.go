package user

import (
	"errors"
	"time"

	"github.com/wellingtonlope/ticket-api/internal/app/repository"
)

var (
	ErrUserEmailPasswordWrong = errors.New("Email or password is wrong")
)

type Login struct {
	userRepository repository.UserRepository
}

func NewLogin(userRepository repository.UserRepository) *Login {
	return &Login{userRepository: userRepository}
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	ID        string
	Name      string
	Email     string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (u *Login) Handle(input LoginInput) (*LoginOutput, error) {
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

	return &LoginOutput{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
