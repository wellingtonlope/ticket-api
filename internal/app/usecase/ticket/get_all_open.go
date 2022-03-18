package ticket

import (
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type GetAllOpen struct {
	ticketRepository repository.TicketRepository
}

func NewGetAllOpen(ticketRepository repository.TicketRepository) *GetAllOpen {
	return &GetAllOpen{ticketRepository: ticketRepository}
}

type GetAllOpenInput struct {
	LoggedUser domain.User
}

func (u *GetAllOpen) Handle(input GetAllOpenInput) (*[]TicketOutput, error) {
	if input.LoggedUser.Profile != domain.PROFILE_OPERATOR {
		return nil, security.ErrForbidden
	}

	tickets, err := u.ticketRepository.GetAllOpen()
	if err != nil {
		return nil, err
	}

	return ticketsOutputsFromTickets(tickets), nil
}
