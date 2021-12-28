package user

import (
	"github.com/wellingtonlope/ticket-api/application/repository"
	"github.com/wellingtonlope/ticket-api/application/usecase/token"
)

type UserUseCase struct {
	UserRepository repository.UserRepository
	TokenUseCase   *token.TokenUseCase
}

type UserReponse struct {
	Token string `json:"token"`
}
