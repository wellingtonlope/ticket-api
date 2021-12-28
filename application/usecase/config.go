package usecase

import (
	"github.com/wellingtonlope/ticket-api/application/repository"
	"github.com/wellingtonlope/ticket-api/application/usecase/ticket"
	"github.com/wellingtonlope/ticket-api/application/usecase/token"
	"github.com/wellingtonlope/ticket-api/application/usecase/user"
	"time"
)

type AllUseCases struct {
	TokenUsecase  *token.TokenUseCase
	UserUseCase   *user.UserUseCase
	TicketUseCase *ticket.TicketUseCase
}

func GetUseCases(repositories repository.Repositories, secret string, duration time.Duration) (*AllUseCases, error) {
	repos, err := repositories.GetRepositories()
	if err != nil {
		return nil, err
	}

	tokenUsecase := &token.TokenUseCase{
		Secret:         secret,
		Duration:       duration,
		UserRepository: repos.UserRepository,
	}

	userUseCase := &user.UserUseCase{
		UserRepository: repos.UserRepository,
		TokenUseCase:   tokenUsecase,
	}

	ticketUseCase := &ticket.TicketUseCase{
		TicketRepository: repos.TicketRepository,
		TokenUseCase:     tokenUsecase,
		UserUseCase:      userUseCase,
	}

	return &AllUseCases{
		TokenUsecase:  tokenUsecase,
		UserUseCase:   userUseCase,
		TicketUseCase: ticketUseCase,
	}, nil
}
