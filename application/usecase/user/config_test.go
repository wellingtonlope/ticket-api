package user

import (
	"time"

	"github.com/wellingtonlope/ticket-api/application/usecase/token"
	"github.com/wellingtonlope/ticket-api/framework/db/local"
)

func newUserUC() *UserUseCase {
	userRepo := &local.UserRepositoryLocal{}
	return &UserUseCase{
		UserRepository: userRepo,
		TokenUseCase: &token.TokenUseCase{
			Secret:         "123",
			Duration:       time.Hour,
			UserRepository: userRepo,
		},
	}
}
