package ticket

import (
	"time"

	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type Close interface {
	Handle(input CloseInput) (*Output, error)
}

type close struct {
	ticketRepository repository.TicketRepository
}

func NewClose(ticketRepository repository.TicketRepository) *close {
	return &close{ticketRepository: ticketRepository}
}

type CloseInput struct {
	TicketID   string
	Solution   string
	UpdatedAt  time.Time
	LoggedUser security.User
}

func (u *close) Handle(input CloseInput) (*Output, error) {
	if input.LoggedUser.Profile != string(domain.ProfileOperator) {
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
