package repository

import (
	"errors"

	"github.com/wellingtonlope/ticket-api/internal/domain"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	Insert(user domain.User) (*domain.User, error)
	GetByID(id string) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	GetAllOperator() (*[]domain.User, error)
}
