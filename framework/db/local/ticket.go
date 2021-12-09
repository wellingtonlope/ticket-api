package local

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
)

type TicketRepositoryLocal struct {
	Tickets []domain.Ticket
}

func (r *TicketRepositoryLocal) Insert(ticket *domain.Ticket) (*domain.Ticket, *myerrors.Error) {
	ticketGot, err := r.GetById(ticket.ID)
	if ticketGot != nil && err == nil {
		return nil, myerrors.NewErrorMessage("ticket already exist", myerrors.REGISTER_ALREADY_EXISTS)
	}
	if err != nil && err.Type != myerrors.REGISTER_NOT_FOUND {
		return nil, myerrors.NewError(err, myerrors.REPOSITORY)
	}

	r.Tickets = append(r.Tickets, *ticket)
	return ticket, nil
}

func (r *TicketRepositoryLocal) GetById(id string) (*domain.Ticket, *myerrors.Error) {
	for _, ticket := range r.Tickets {
		if ticket.ID == id {
			return &ticket, nil
		}
	}

	return nil, myerrors.NewErrorMessage("ticket not found", myerrors.REGISTER_NOT_FOUND)
}

func (r *TicketRepositoryLocal) Update(ticket *domain.Ticket) (*domain.Ticket, *myerrors.Error) {
	ticketGot, err := r.GetById(ticket.ID)
	if ticketGot == nil && err != nil && err.Type == myerrors.REGISTER_NOT_FOUND {
		return nil, myerrors.NewErrorMessage("ticket not found", myerrors.REGISTER_NOT_FOUND)
	}
	if err != nil && err.Type != myerrors.REGISTER_NOT_FOUND {
		return nil, myerrors.NewError(err, myerrors.REPOSITORY)
	}

	for index, item := range r.Tickets {
		if item.ID == ticket.ID {
			r.Tickets[index] = *ticket
			break
		}
	}

	return ticket, nil
}

func (r *TicketRepositoryLocal) GetAll() (*[]domain.Ticket, *myerrors.Error) {
	return &r.Tickets, nil
}

func (r *TicketRepositoryLocal) GetAllOpen() (*[]domain.Ticket, *myerrors.Error) {
	filtered := []domain.Ticket{}

	for _, ticket := range r.Tickets {
		if ticket.Status == domain.STATUS_OPEN {
			filtered = append(filtered, ticket)
		}
	}

	return &filtered, nil
}

func (r *TicketRepositoryLocal) GetAllByOperator(operator *domain.User) (*[]domain.Ticket, *myerrors.Error) {
	filtered := []domain.Ticket{}

	for _, ticket := range r.Tickets {
		if ticket.Operator != nil && ticket.Operator.ID == operator.ID {
			filtered = append(filtered, ticket)
		}
	}

	return &filtered, nil
}

func (r *TicketRepositoryLocal) GetAllByClient(client *domain.User) (*[]domain.Ticket, *myerrors.Error) {
	filtered := []domain.Ticket{}

	for _, ticket := range r.Tickets {
		if ticket.Client != nil && ticket.Client.ID == client.ID {
			filtered = append(filtered, ticket)
		}
	}

	return &filtered, nil
}

func (r *TicketRepositoryLocal) Delete(id string) *myerrors.Error {
	_, err := r.GetById(id)
	if err != nil {
		return myerrors.NewErrorMessage("ticket not found", myerrors.REGISTER_NOT_FOUND)
	}

	for index, ticket := range r.Tickets {
		if ticket.ID == id {
			r.Tickets = removeIndexTicket(r.Tickets, index)
		}
	}

	return nil
}

func removeIndexTicket(s []domain.Ticket, index int) []domain.Ticket {
	return append(s[:index], s[index+1:]...)
}
