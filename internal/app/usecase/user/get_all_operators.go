package user

import (
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type GetAllOperators interface {
	Handle(input GetAllOperatorsInput) (*[]Output, error)
}

type getAllOperators struct {
	userRepository repository.UserRepository
}

func NewGetAllOperators(userRepository repository.UserRepository) GetAllOperators {
	return &getAllOperators{userRepository: userRepository}
}

type GetAllOperatorsInput struct {
	LoggedUser security.User
}

func (u *getAllOperators) Handle(input GetAllOperatorsInput) (*[]Output, error) {
	if input.LoggedUser.Profile != string(domain.ProfileOperator) {
		return nil, security.ErrForbidden
	}

	users, err := u.userRepository.GetAllOperator()
	if err != nil {
		return nil, err
	}

	return usersOutputsFromUsers(users), nil
}
