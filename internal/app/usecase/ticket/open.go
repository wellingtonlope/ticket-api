package ticket

import (
	"time"

	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type Open struct {
	ticketRepository repository.TicketRepository
}

func NewOpen(ticketRepository repository.TicketRepository) *Open {
	return &Open{ticketRepository: ticketRepository}
}

type OpenInput struct {
	Title       string
	Description string
	CreatedAt   time.Time
	LoggedUser  domain.User
}

func (u *Open) Handle(input OpenInput) (*TicketOutput, error) {
	loggedUser := input.LoggedUser
	ticket, err := domain.OpenTicket(input.Title, input.Description, input.CreatedAt, loggedUser)
	if err != nil {
		return nil, err
	}

	ticket, err = u.ticketRepository.Insert(*ticket)
	if err != nil {
		return nil, err
	}

	return ticketOutputFromTicket(ticket), nil
}
