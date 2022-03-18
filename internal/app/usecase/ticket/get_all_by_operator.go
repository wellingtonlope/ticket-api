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
	LoggedUser domain.User
}

func (u *GetAllByOperator) Handle(input GetAllByOperatorInput) (*[]TicketOutput, error) {
	if input.LoggedUser.Profile != domain.PROFILE_OPERATOR {
		return nil, security.ErrForbidden
	}

	tickets, err := u.ticketRepository.GetAllByOperatorID(input.OperatorID)
	if err != nil {
		return nil, err
	}
	return ticketsOutputsFromTickets(tickets), nil
}
