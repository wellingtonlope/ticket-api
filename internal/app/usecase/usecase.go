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

func GetUseCases(repositories *repository.AllRepositories) (*AllUseCases, error) {
	return &AllUseCases{
		UserGetAllOperators:    user.NewGetAllOperators(repositories.UserRepository),
		UserLogin:              user.NewLogin(repositories.UserRepository),
		UserRegister:           user.NewRegister(repositories.UserRepository),
		TicketAssignToOperator: ticket.NewAssignToOperator(repositories.TicketRepository, repositories.UserRepository),
		TicketClose:            ticket.NewClose(repositories.TicketRepository),
		TicketDelete:           ticket.NewDelete(repositories.TicketRepository),
		TicketGet:              ticket.NewGet(repositories.TicketRepository),
		TicketGetAll:           ticket.NewGetAll(repositories.TicketRepository),
		TicketGetAllByClient:   ticket.NewGetAllByClient(repositories.TicketRepository),
		TicketGetAllByOperator: ticket.NewGetAllByOperator(repositories.TicketRepository),
		TicketGetAllOPen:       ticket.NewGetAllOpen(repositories.TicketRepository),
		TicketGetByID:          ticket.NewGetByID(repositories.TicketRepository),
		TicketOpen:             ticket.NewOpen(repositories.TicketRepository),
	}, nil
}
