package ticket

import (
	"time"

	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type Get struct {
	ticketRepository repository.TicketRepository
}

func NewGet(ticketRepository repository.TicketRepository) *Get {
	return &Get{ticketRepository: ticketRepository}
}

type GetInput struct {
	TicketID   string
	UpdatedAt  time.Time
	LoggedUser domain.User
}

func (u *Get) Handle(input GetInput) (*TicketOutput, error) {
	ticket, err := u.ticketRepository.GetByID(input.TicketID)
	if err != nil {
		return nil, err
	}

	err = ticket.Get(input.LoggedUser, input.UpdatedAt)
	if err != nil {
		if err == domain.ErrTicketNoOperator {
			return nil, security.ErrForbidden
		}
		return nil, err
	}

	ticket, err = u.ticketRepository.Update(*ticket)
	if err != nil {
		return nil, err
	}

	return ticketOutputFromTicket(ticket), nil
}
