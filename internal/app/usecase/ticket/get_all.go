package ticket

import (
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type GetAll struct {
	ticketRepository repository.TicketRepository
}

func NewGetAll(ticketRepository repository.TicketRepository) *GetAll {
	return &GetAll{ticketRepository: ticketRepository}
}

type GetAllInput struct {
	LoggedUser security.User
}

func (u *GetAll) Handle(input GetAllInput) (*[]Output, error) {
	if input.LoggedUser.Profile != string(domain.ProfileOperator) {
		return nil, security.ErrForbidden
	}

	tickets, err := u.ticketRepository.GetAll()
	if err != nil {
		return nil, err
	}

	return ticketsOutputsFromTickets(tickets), nil
}
