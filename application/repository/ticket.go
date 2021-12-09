package repository

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
)

type TicketRepository interface {
	GetById(id string) (*domain.Ticket, *myerrors.Error)
	GetAll() (*[]domain.Ticket, *myerrors.Error)
	GetAllOpen() (*[]domain.Ticket, *myerrors.Error)
	GetAllByOperator(operator *domain.User) (*[]domain.Ticket, *myerrors.Error)
	GetAllByClient(client *domain.User) (*[]domain.Ticket, *myerrors.Error)
	Insert(ticket *domain.Ticket) (*domain.Ticket, *myerrors.Error)
	Update(ticket *domain.Ticket) (*domain.Ticket, *myerrors.Error)
	Delete(id string) *myerrors.Error
}
