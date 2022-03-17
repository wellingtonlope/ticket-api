package user

import (
	"time"

	"github.com/wellingtonlope/ticket-api/internal/app"
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type GetAllOperators struct {
	userRepository repository.UserRepository
}

func NewGetAllOperators(userRepository repository.UserRepository) *GetAllOperators {
	return &GetAllOperators{userRepository: userRepository}
}

type GetAllOperatorsInput struct {
	LoggedUser domain.User
}

type GetAllOperatorsOutput struct {
	ID        string
	Name      string
	Email     string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (u *GetAllOperators) Handle(input GetAllOperatorsInput) (*[]GetAllOperatorsOutput, error) {
	if input.LoggedUser.Profile != domain.PROFILE_OPERATOR {
		return nil, app.ErrForbidden
	}

	users, err := u.userRepository.GetAllOperator()
	if err != nil {
		return nil, err
	}

	var output []GetAllOperatorsOutput
	for _, user := range *users {
		output = append(output, GetAllOperatorsOutput{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email.String(),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return &output, nil
}
