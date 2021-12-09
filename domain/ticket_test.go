package domain

import "testing"

func TestTicketOpen(t *testing.T) {
	title, description := "title", "description"
	client, _ := UserRegister("client", "client@mail.com", "password")

	t.Run("a valid ticket", func(t *testing.T) {
		got, _ := TicketOpen(title, description, client)

		if got.ID == "" {
			t.Error("expected an ID, but got nothing")
		}

		if got.CreatedAt == nil {
			t.Error("expected a created date, but got a nil")
		}

		if got.UpdatedAt != nil {
			t.Error("expected a nil, but got a date")
		}

		if got.Title != title {
			t.Errorf("expected %q, but got %q", title, got.Title)
		}

		if got.Description != description {
			t.Errorf("expected %q, but got %q", title, got.Title)
		}

		if got.Solution != "" {
			t.Errorf("expected an empty solution, but got %q", got.Solution)
		}

		if got.Status != STATUS_OPEN {
			t.Errorf("expected %q, but got %q", STATUS_OPEN, got.Status)
		}

		if got.Client.ID != client.ID || got.Client.Name != client.Name || got.Client.Email != client.Email {
			t.Errorf("expected %v, but got %v", client, got.Client)
		}
	})

	t.Run("empty title", func(t *testing.T) {
		_, err := TicketOpen("", description, client)
		if err == nil {
			t.Error("expected an error, but got nil")
		}
	})

	t.Run("nil client", func(t *testing.T) {
		_, err := TicketOpen(title, description, nil)
		if err == nil {
			t.Error("expected an error, but got nil")
		}
	})

	t.Run("client is not a client", func(t *testing.T) {
		operator, _ := UserRegister("operator", "operator@mail.com", "password")
		operator.Profile = PROFILE_OPERATOR
		_, err := TicketOpen(title, description, operator)
		if err == nil {
			t.Error("expected an error, but got nil")
		}
	})
}

func TestTicketGet(t *testing.T) {
	title, description := "title", "description"
	client, _ := UserRegister("client", "client@mail.com", "password")
	operator, _ := UserRegister("operator", "operator@mail.com", "password")
	operator.Profile = PROFILE_OPERATOR

	t.Run("a valid get a ticket", func(t *testing.T) {
		ticket, _ := TicketOpen(title, description, client)
		ticket.Get(operator)

		if ticket.Status != STATUS_IN_PROGRESS {
			t.Errorf("expected %q, but got %q", STATUS_IN_PROGRESS, ticket.Status)
		}

		if ticket.UpdatedAt == nil {
			t.Error("expected a date, but got nil")
		}
	})

	t.Run("nil operator", func(t *testing.T) {
		ticket, _ := TicketOpen(title, description, client)
		err := ticket.Get(nil)
		if err == nil {
			t.Error("expected an error, but got a nil")
		}
	})

	t.Run("operator is not an operator", func(t *testing.T) {
		ticket, _ := TicketOpen(title, description, client)
		err := ticket.Get(client)
		if err == nil {
			t.Error("expected an error, but got a nil")
		}
	})
}

func TestTicketClose(t *testing.T) {
	title, description := "title", "description"
	client, _ := UserRegister("client", "client@mail.com", "password")
	operator, _ := UserRegister("operator", "operator@mail.com", "password")
	operator.Profile = PROFILE_OPERATOR

	t.Run("a valid close a ticket", func(t *testing.T) {
		ticket, _ := TicketOpen(title, description, client)
		ticket.Get(operator)
		ticket.Close("solution")

		if ticket.Status != STATUS_CLOSE {
			t.Errorf("expected %q, but got %q", STATUS_CLOSE, ticket.Status)
		}

		if ticket.UpdatedAt == nil {
			t.Error("expected a date, but got nil")
		}

	})

	t.Run("a valid close a ticket", func(t *testing.T) {
		ticket, _ := TicketOpen(title, description, client)
		err := ticket.Close("solution")
		if err == nil {
			t.Error("expected an error, but got a nil")
		}
	})

}
