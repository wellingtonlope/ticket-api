package ticket

import (
	"time"

	"github.com/wellingtonlope/ticket-api/application/usecase/token"
	"github.com/wellingtonlope/ticket-api/application/usecase/user"
	"github.com/wellingtonlope/ticket-api/framework/db/local"
)

func newTicketUC() *TicketUseCase {
	userRepository := &local.UserRepositoryLocal{}
	tokenUseCase := &token.TokenUseCase{
		Secret:         "123",
		Duration:       time.Hour,
		UserRepository: userRepository,
	}
	userUseCase := &user.UserUseCase{
		UserRepository: userRepository,
		TokenUseCase:   tokenUseCase,
	}

	return &TicketUseCase{
		TicketRepository: &local.TicketRepositoryLocal{},
		UserUseCase:      userUseCase,
		TokenUseCase:     tokenUseCase,
	}
}
