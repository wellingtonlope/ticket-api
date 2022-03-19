package user

import (
	"time"

	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type UserOutput struct {
	ID        string
	Name      string
	Email     string
	Profile   string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func userOutputFromUser(user *domain.User) *UserOutput {
	return &UserOutput{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email.String(),
		Profile:   string(user.Profile),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func usersOutputsFromUsers(users *[]domain.User) *[]UserOutput {
	outputs := make([]UserOutput, 0, len(*users))
	for _, user := range *users {
		outputs = append(outputs, *userOutputFromUser(&user))
	}
	return &outputs
}
