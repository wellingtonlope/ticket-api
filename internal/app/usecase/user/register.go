package user

import (
	"errors"
	"time"

	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type Register struct {
	userRepository repository.UserRepository
}

func NewRegister(userRepository repository.UserRepository) *Register {
	return &Register{userRepository: userRepository}
}

type RegisterInput struct {
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

func (u *Register) Handle(input RegisterInput) (*Output, error) {
	user, err := u.userRepository.GetByEmail(input.Email)
	if err != nil && err != repository.ErrUserNotFound {
		return nil, err
	}

	if user != nil {
		return nil, ErrUserAlreadyExists
	}

	user, err = domain.UserRegister(input.Name, input.Email, input.Password, input.CreatedAt)
	if err != nil {
		return nil, err
	}

	user, err = u.userRepository.Insert(*user)
	if err != nil {
		return nil, err
	}

	return userOutputFromUser(user), nil
}
