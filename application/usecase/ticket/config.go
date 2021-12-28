package ticket

import (
	"github.com/wellingtonlope/ticket-api/application/repository"
	"github.com/wellingtonlope/ticket-api/application/usecase/token"
	"github.com/wellingtonlope/ticket-api/application/usecase/user"
)

type TicketUseCase struct {
	TicketRepository repository.TicketRepository
	UserUseCase      *user.UserUseCase
	TokenUseCase     *token.TokenUseCase
}
