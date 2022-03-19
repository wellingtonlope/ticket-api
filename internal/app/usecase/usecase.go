package usecase

import (
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/app/usecase/ticket"
	"github.com/wellingtonlope/ticket-api/internal/app/usecase/user"
)

type AllUseCases struct {
	UserGetAllOperators    *user.GetAllOperators
	UserLogin              *user.Login
	UserRegister           *user.Register
	TicketAssignToOperator *ticket.AssignToOperator
	TicketClose            *ticket.Close
	TicketDelete           *ticket.Delete
	TicketGet              *ticket.Get
	TicketGetAll           *ticket.GetAll
	TicketGetAllByClient   *ticket.GetAllByClient
	TicketGetAllByOperator *ticket.GetAllByOperator
	TicketGetAllOPen       *ticket.GetAllOpen
	TicketGetByID          *ticket.GetByID
	TicketOpen             *ticket.Open
}

func GetUseCases(repositories repository.Repositories) (*AllUseCases, error) {
	repos, err := repositories.GetRepositories()
	if err != nil {
		return nil, err
	}

	return &AllUseCases{
		UserGetAllOperators:    user.NewGetAllOperators(repos.UserRepository),
		UserLogin:              user.NewLogin(repos.UserRepository),
		UserRegister:           user.NewRegister(repos.UserRepository),
		TicketAssignToOperator: ticket.NewAssignToOperator(repos.TicketRepository, repos.UserRepository),
		TicketClose:            ticket.NewClose(repos.TicketRepository),
		TicketDelete:           ticket.NewDelete(repos.TicketRepository),
		TicketGet:              ticket.NewGet(repos.TicketRepository),
		TicketGetAll:           ticket.NewGetAll(repos.TicketRepository),
		TicketGetAllByClient:   ticket.NewGetAllByClient(repos.TicketRepository),
		TicketGetAllByOperator: ticket.NewGetAllByOperator(repos.TicketRepository),
		TicketGetAllOPen:       ticket.NewGetAllOpen(repos.TicketRepository),
		TicketGetByID:          ticket.NewGetByID(repos.TicketRepository),
		TicketOpen:             ticket.NewOpen(repos.TicketRepository),
	}, nil
}
