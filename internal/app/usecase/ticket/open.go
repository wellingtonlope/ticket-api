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

type OpenUserOutput struct {
	ID    string
	Name  string
	Email string
}

type OpenOutput struct {
	ID          string
	Title       string
	Description string
	Solution    string
	Status      string
	Client      *OpenUserOutput
	CreatedAt   *time.Time
}

func (u *Open) Handle(input OpenInput) (*OpenOutput, error) {
	loggedUser := input.LoggedUser
	ticket, err := domain.OpenTicket(input.Title, input.Description, input.CreatedAt, loggedUser)
	if err != nil {
		return nil, err
	}

	ticket, err = u.ticketRepository.Insert(*ticket)
	if err != nil {
		return nil, err
	}

	return &OpenOutput{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Solution:    ticket.Solution,
		Status:      string(ticket.Status),
		Client: &OpenUserOutput{
			ID:    loggedUser.ID,
			Name:  loggedUser.Name,
			Email: loggedUser.Email.String(),
		},
		CreatedAt: ticket.CreatedAt,
	}, nil
}
