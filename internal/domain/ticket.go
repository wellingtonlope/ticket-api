package domain

import (
	"errors"
	"time"
)

type Status string

const (
	StatusOpen       Status = "OPEN"
	StatusInProgress Status = "IN_PROGRESS"
	StatusClose      Status = "CLOSE"
)

var (
	ErrTicketNoOperator     = errors.New("operator must be an operator")
	ErrTicketTitleIsInvalid = errors.New("title mustn't be empty")
	ErrTicketNoGetToClose   = errors.New("first you need to get a ticket")
)

type TicketUser struct {
	ID    string
	Name  string
	Email Email
}

type Ticket struct {
	ID          string
	Title       string
	Description string
	Solution    string
	Status      Status
	Client      TicketUser
	Operator    *TicketUser
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func OpenTicket(title, description string, createdAt time.Time, client User) (Ticket, error) {
	if title == "" {
		return Ticket{}, ErrTicketTitleIsInvalid
	}

	return Ticket{
		Title:       title,
		Description: description,
		Status:      StatusOpen,
		Client: TicketUser{
			ID:    client.ID,
			Name:  client.Name,
			Email: client.Email,
		},
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}, nil
}

func (t *Ticket) Get(operator User, updatedAt time.Time) error {
	if operator.Profile != ProfileOperator {
		return ErrTicketNoOperator
	}

	t.Operator = &TicketUser{
		ID:    operator.ID,
		Name:  operator.Name,
		Email: operator.Email,
	}
	t.UpdatedAt = updatedAt
	t.Status = StatusInProgress

	return nil
}

func (t *Ticket) Close(solution string, updatedAt time.Time) error {
	if t.Operator == nil {
		return ErrTicketNoGetToClose
	}

	t.Solution = solution
	t.UpdatedAt = updatedAt
	t.Status = StatusClose

	return nil
}
