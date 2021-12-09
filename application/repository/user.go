package repository

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
)

type UserRepository interface {
	Insert(ticket *domain.User) (*domain.User, *myerrors.Error)
	GetById(id string) (*domain.User, *myerrors.Error)
	GetByEmail(email string) (*domain.User, *myerrors.Error)
	GetAllOperator() (*[]domain.User, *myerrors.Error)
	Delete(id string) *myerrors.Error
}
