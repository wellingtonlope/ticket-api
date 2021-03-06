package ticket

import (
	"time"

	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type Get interface {
	Handle(input GetInput) (*Output, error)
}

type get struct {
	ticketRepository repository.TicketRepository
	userRepository   repository.UserRepository
}

func NewGet(ticketRepository repository.TicketRepository, userRepository repository.UserRepository) Get {
	return &get{ticketRepository: ticketRepository, userRepository: userRepository}
}

type GetInput struct {
	TicketID   string
	UpdatedAt  time.Time
	LoggedUser security.User
}

func (u *get) Handle(input GetInput) (*Output, error) {
	ticket, err := u.ticketRepository.GetByID(input.TicketID)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepository.GetByID(input.LoggedUser.ID)
	if err != nil {
		return nil, err
	}

	err = ticket.Get(*user, input.UpdatedAt)
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
