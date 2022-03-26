package ticket

import (
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type GetAllByClient interface {
	Handle(input GetAllByClientInput) (*[]Output, error)
}

type getAllByClient struct {
	ticketRepository repository.TicketRepository
}

func NewGetAllByClient(ticketRepository repository.TicketRepository) GetAllByClient {
	return &getAllByClient{ticketRepository: ticketRepository}
}

type GetAllByClientInput struct {
	ClientID   string
	LoggedUser security.User
}

func (u *getAllByClient) Handle(input GetAllByClientInput) (*[]Output, error) {
	if input.ClientID != input.LoggedUser.ID && input.LoggedUser.Profile != string(domain.ProfileOperator) {
		return nil, security.ErrForbidden
	}

	tickets, err := u.ticketRepository.GetAllByClientID(input.ClientID)
	if err != nil {
		return nil, err
	}

	return ticketsOutputsFromTickets(tickets), nil
}
