package usecase

import (
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/app/usecase/ticket"
	"github.com/wellingtonlope/ticket-api/internal/app/usecase/user"
)

type AllUseCases struct {
	userGetAllOperators    *user.GetAllOperators
	userLogin              *user.Login
	userRegister           *user.Register
	ticketAssignToOperator *ticket.AssignToOperator
	ticketClose            *ticket.Close
	ticketDelete           *ticket.Delete
	ticketGet              *ticket.Get
	ticketGetAll           *ticket.GetAll
	ticketGetAllByClient   *ticket.GetAllByClient
	ticketGetAllByOperator *ticket.GetAllByOperator
	ticketGetAllOPen       *ticket.GetAllOpen
	ticketGetByID          *ticket.GetByID
	ticketOpen             *ticket.Open
}

func GetUseCases(repositories repository.Repositories) (*AllUseCases, error) {
	repos, err := repositories.GetRepositories()
	if err != nil {
		return nil, err
	}

	return &AllUseCases{
		userGetAllOperators:    user.NewGetAllOperators(repos.UserRepository),
		userLogin:              user.NewLogin(repos.UserRepository),
		userRegister:           user.NewRegister(repos.UserRepository),
		ticketAssignToOperator: ticket.NewAssignToOperator(repos.TicketRepository, repos.UserRepository),
		ticketClose:            ticket.NewClose(repos.TicketRepository),
		ticketDelete:           ticket.NewDelete(repos.TicketRepository),
		ticketGet:              ticket.NewGet(repos.TicketRepository),
		ticketGetAll:           ticket.NewGetAll(repos.TicketRepository),
		ticketGetAllByClient:   ticket.NewGetAllByClient(repos.TicketRepository),
		ticketGetAllByOperator: ticket.NewGetAllByOperator(repos.TicketRepository),
		ticketGetAllOPen:       ticket.NewGetAllOpen(repos.TicketRepository),
		ticketGetByID:          ticket.NewGetByID(repos.TicketRepository),
		ticketOpen:             ticket.NewOpen(repos.TicketRepository),
	}, nil
}
