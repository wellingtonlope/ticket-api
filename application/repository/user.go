package repository

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
)

type UserRepository interface {
	GetById(id string) (*domain.User, *myerrors.Error)
	GetByEmail(email string) (*domain.User, *myerrors.Error)
	GetAllOperators() (*[]domain.User, *myerrors.Error)
}
