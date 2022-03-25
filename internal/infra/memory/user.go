package memory

import (
	"github.com/google/uuid"
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type UserRepository struct {
	users []domain.User
}

func (r *UserRepository) Insert(user domain.User) (*domain.User, error) {
	user.ID = uuid.New().String()

	r.users = append(r.users, user)
	return &user, nil
}

func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, repository.ErrUserNotFound
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	for _, user := range r.users {
		if user.Email.String() == email {
			return &user, nil
		}
	}

	return nil, repository.ErrUserNotFound
}

func (r *UserRepository) GetAllOperator() (*[]domain.User, error) {
	operators := make([]domain.User, 0, 10)

	for _, user := range r.users {
		if user.Profile == domain.ProfileOperator {
			operators = append(operators, user)
		}
	}

	return &operators, nil
}
