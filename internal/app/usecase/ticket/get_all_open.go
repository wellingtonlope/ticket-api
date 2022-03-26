package ticket

import (
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type GetAllOpen interface {
	Handle(input GetAllOpenInput) (*[]Output, error)
}

type getAllOpen struct {
	ticketRepository repository.TicketRepository
}

func NewGetAllOpen(ticketRepository repository.TicketRepository) GetAllOpen {
	return &getAllOpen{ticketRepository: ticketRepository}
}

type GetAllOpenInput struct {
	LoggedUser security.User
}

func (u *getAllOpen) Handle(input GetAllOpenInput) (*[]Output, error) {
	if input.LoggedUser.Profile != string(domain.ProfileOperator) {
		return nil, security.ErrForbidden
	}

	tickets, err := u.ticketRepository.GetAllOpen()
	if err != nil {
		return nil, err
	}

	return ticketsOutputsFromTickets(tickets), nil
}
