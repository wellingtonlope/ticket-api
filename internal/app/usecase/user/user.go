package user

import (
	"time"

	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type Output struct {
	ID        string
	Name      string
	Email     string
	Profile   string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func userOutputFromUser(user *domain.User) *Output {
	return &Output{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email.String(),
		Profile:   string(user.Profile),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func usersOutputsFromUsers(users *[]domain.User) *[]Output {
	outputs := make([]Output, 0, len(*users))
	for _, user := range *users {
		outputs = append(outputs, *userOutputFromUser(&user))
	}
	return &outputs
}
