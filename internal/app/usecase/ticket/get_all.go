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
	LoggedUser domain.User
}

func (u *GetAll) Handle(input GetAllInput) (*[]TicketOutput, error) {
	if input.LoggedUser.Profile != domain.PROFILE_OPERATOR {
		return nil, security.ErrForbidden
	}

	tickets, err := u.ticketRepository.GetAll()
	if err != nil {
		return nil, err
	}

	return ticketsOutputsFromTickets(tickets), nil
}
