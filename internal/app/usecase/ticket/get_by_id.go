package ticket

import (
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type GetByID interface {
	Handle(input GetByIDInput) (*Output, error)
}

type getByID struct {
	ticketRepository repository.TicketRepository
}

func NewGetByID(ticketRepository repository.TicketRepository) GetByID {
	return &getByID{ticketRepository: ticketRepository}
}

type GetByIDInput struct {
	TicketID   string
	LoggedUser security.User
}

func (u *getByID) Handle(input GetByIDInput) (*Output, error) {
	loggedUser := input.LoggedUser
	ticket, err := u.ticketRepository.GetByID(input.TicketID)
	if err != nil {
		return nil, err
	}

	if ticket.Client.ID != loggedUser.ID && loggedUser.Profile != string(domain.ProfileOperator) {
		return nil, repository.ErrTicketNotFound
	}

	return ticketOutputFromTicket(ticket), nil
}
