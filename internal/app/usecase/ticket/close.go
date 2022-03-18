package ticket

import (
	"time"

	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type Close struct {
	ticketRepository repository.TicketRepository
}

func NewClose(ticketRepository repository.TicketRepository) *Close {
	return &Close{ticketRepository: ticketRepository}
}

type CloseInput struct {
	TicketID   string
	Solution   string
	UpdatedAt  time.Time
	LoggedUser domain.User
}

func (u *Close) Handle(input CloseInput) (*TicketOutput, error) {
	if input.LoggedUser.Profile != domain.PROFILE_OPERATOR {
		return nil, security.ErrForbidden
	}

	ticket, err := u.ticketRepository.GetByID(input.TicketID)
	if err != nil {
		return nil, err
	}

	err = ticket.Close(input.Solution, input.UpdatedAt)
	if err != nil {
		return nil, err
	}

	ticket, err = u.ticketRepository.Update(*ticket)
	if err != nil {
		return nil, err
	}

	return ticketOutputFromTicket(ticket), nil
}
