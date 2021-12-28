package user

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
)

func (uc *UserUseCase) GetById(id string) (*domain.User, *myerrors.Error) {
	return uc.UserRepository.GetById(id)
}
