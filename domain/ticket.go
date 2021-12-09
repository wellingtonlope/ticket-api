package domain

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Status string

const (
	STATUS_OPEN        Status = "OPEN"
	STATUS_IN_PROGRESS Status = "IN_PROGRESS"
	STATUS_CLOSE       Status = "CLOSE"
)

type TicketUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Ticket struct {
	Base        `bson:",inline"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Solution    string      `json:"solution"`
	Status      Status      `json:"status"`
	Client      *TicketUser `json:"client"`
	Operator    *TicketUser `json:"operator"`
}

func TicketOpen(title, description string, client *User) (*Ticket, error) {
	if title == "" {
		return nil, errors.New("title musn't be empty")
	}

	if client == nil {
		return nil, errors.New("client musn't be nil")
	}

	if client.Profile != PROFILE_CLIENT {
		return nil, errors.New("client must be a client")
	}

	now := time.Now()
	return &Ticket{
		Base: Base{
			ID:        uuid.NewV4().String(),
			CreatedAt: &now,
		},
		Title:       title,
		Description: description,
		Status:      STATUS_OPEN,
		Client: &TicketUser{
			ID:    client.ID,
			Name:  client.Name,
			Email: client.Email,
		},
	}, nil
}

func (t *Ticket) Get(operator *User) error {
	if operator == nil {
		return errors.New("operator musn't be nil")
	}

	if operator.Profile != PROFILE_OPERATOR {
		return errors.New("operator must be an operator")
	}

	t.Operator = &TicketUser{
		ID:    operator.ID,
		Name:  operator.Name,
		Email: operator.Email,
	}
	now := time.Now()
	t.UpdatedAt = &now

	return nil
}

func (t *Ticket) Close(solution string) error {
	if t.Operator == nil {
		return errors.New("first you need to get a ticket")
	}

	t.Solution = solution
	now := time.Now()
	t.UpdatedAt = &now

	return nil
}
