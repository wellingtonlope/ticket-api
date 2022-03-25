package domain

import (
	"errors"
	"time"
)

type Status string

const (
	STATUS_OPEN        Status = "OPEN"
	STATUS_IN_PROGRESS Status = "IN_PROGRESS"
	STATUS_CLOSE       Status = "CLOSE"
)

var (
	ErrTicketNoOperator     = errors.New("operator must be an operator")
	ErrTicketTitleIsInvalid = errors.New("title musn't be empty")
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
	Client      *TicketUser
	Operator    *TicketUser
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

func OpenTicket(title, description string, createdAt time.Time, client User) (*Ticket, error) {
	if title == "" {
		return nil, ErrTicketTitleIsInvalid
	}

	return &Ticket{
		Title:       title,
		Description: description,
		Status:      STATUS_OPEN,
		Client: &TicketUser{
			ID:    client.ID,
			Name:  client.Name,
			Email: client.Email,
		},
		CreatedAt: &createdAt,
	}, nil
}

func (t *Ticket) Get(operator User, updatedAt time.Time) error {
	if operator.Profile != PROFILE_OPERATOR {
		return ErrTicketNoOperator
	}

	t.Operator = &TicketUser{
		ID:    operator.ID,
		Name:  operator.Name,
		Email: operator.Email,
	}
	t.UpdatedAt = &updatedAt
	t.Status = STATUS_IN_PROGRESS

	return nil
}

func (t *Ticket) Close(solution string, updatedAt time.Time) error {
	if t.Operator == nil {
		return ErrTicketNoGetToClose
	}

	t.Solution = solution
	t.UpdatedAt = &updatedAt
	t.Status = STATUS_CLOSE

	return nil
}
