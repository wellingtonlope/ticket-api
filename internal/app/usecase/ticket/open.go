package ticket

import (
	"time"

	"github.com/wellingtonlope/ticket-api/internal/app/security"

	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type Open interface {
	Handle(input OpenInput) (*Output, error)
}

type open struct {
	ticketRepository repository.TicketRepository
	userRepository   repository.UserRepository
}

func NewOpen(ticketRepository repository.TicketRepository, userRepository repository.UserRepository) Open {
	return &open{ticketRepository: ticketRepository, userRepository: userRepository}
}

type OpenInput struct {
	Title       string
	Description string
	CreatedAt   time.Time
	LoggedUser  security.User
}

func (u *open) Handle(input OpenInput) (*Output, error) {
	user, err := u.userRepository.GetByID(input.LoggedUser.ID)
	if err != nil {
		return nil, err
	}

	ticket, err := domain.OpenTicket(input.Title, input.Description, input.CreatedAt, *user)
	if err != nil {
		return nil, err
	}

	ticket, err = u.ticketRepository.Insert(*ticket)
	if err != nil {
		return nil, err
	}

	return ticketOutputFromTicket(ticket), nil
}
