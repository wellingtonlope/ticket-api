package ticket

import (
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type GetAllByOperator struct {
	ticketRepository repository.TicketRepository
}

func NewGetAllByOperator(ticketRepository repository.TicketRepository) *GetAllByOperator {
	return &GetAllByOperator{ticketRepository: ticketRepository}
}

type GetAllByOperatorInput struct {
	OperatorID string
	LoggedUser security.User
}

func (u *GetAllByOperator) Handle(input GetAllByOperatorInput) (*[]Output, error) {
	if input.LoggedUser.Profile != string(domain.ProfileOperator) {
		return nil, security.ErrForbidden
	}

	tickets, err := u.ticketRepository.GetAllByOperatorID(input.OperatorID)
	if err != nil {
		return nil, err
	}
	return ticketsOutputsFromTickets(tickets), nil
}
