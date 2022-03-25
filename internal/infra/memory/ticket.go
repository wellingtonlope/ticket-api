package memory

import (
	"github.com/google/uuid"
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type TicketRepository struct {
	tickets []domain.Ticket
}

func (r *TicketRepository) GetByID(id string) (*domain.Ticket, error) {
	for _, ticket := range r.tickets {
		if ticket.ID == id {
			return &ticket, nil
		}
	}

	return nil, repository.ErrTicketNotFound
}

func (r *TicketRepository) Insert(ticket domain.Ticket) (*domain.Ticket, error) {
	ticket.ID = uuid.New().String()

	r.tickets = append(r.tickets, ticket)
	return &ticket, nil
}

func (r *TicketRepository) Update(ticket domain.Ticket) (*domain.Ticket, error) {
	for i, t := range r.tickets {
		if t.ID == ticket.ID {
			r.tickets[i] = ticket
			return &ticket, nil
		}
	}

	return nil, repository.ErrTicketNotFound
}

func (r *TicketRepository) GetAll() (*[]domain.Ticket, error) {
	return &r.tickets, nil
}

func (r *TicketRepository) GetAllOpen() (*[]domain.Ticket, error) {
	var ticketsOpen []domain.Ticket

	for _, ticket := range r.tickets {
		if ticket.Status == domain.StatusOpen {
			ticketsOpen = append(ticketsOpen, ticket)
		}
	}

	return &ticketsOpen, nil
}

func (r *TicketRepository) GetAllByOperatorID(operatorID string) (*[]domain.Ticket, error) {
	var tickets []domain.Ticket

	for _, ticket := range r.tickets {
		if ticket.Operator != nil && ticket.Operator.ID == operatorID {
			tickets = append(tickets, ticket)
		}
	}

	return &tickets, nil
}

func (r *TicketRepository) GetAllByClientID(clientID string) (*[]domain.Ticket, error) {
	var tickets []domain.Ticket

	for _, ticket := range r.tickets {
		if ticket.Client.ID == clientID {
			tickets = append(tickets, ticket)
		}
	}

	return &tickets, nil
}

func (r *TicketRepository) DeleteByID(id string) error {
	for i, ticket := range r.tickets {
		if ticket.ID == id {
			r.tickets = append(r.tickets[:i], r.tickets[i+1:]...)
			return nil
		}
	}

	return repository.ErrTicketNotFound
}
