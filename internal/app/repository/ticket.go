package repository

import (
	"errors"

	"github.com/wellingtonlope/ticket-api/internal/domain"
)

var (
	ErrTicketNotFound = errors.New("ticket not found")
)

type TicketRepository interface {
	Insert(ticket domain.Ticket) (*domain.Ticket, error)
	GetByID(id string) (*domain.Ticket, error)
	Update(ticket domain.Ticket) (*domain.Ticket, error)
	GetAll() (*[]domain.Ticket, error)
	GetAllOpen() (*[]domain.Ticket, error)
	GetAllByOperatorID(operatorID string) (*[]domain.Ticket, error)
	GetAllByClientID(clientID string) (*[]domain.Ticket, error)
	DeleteByID(id string) error
}
