package local

import (
	"testing"

	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
)

func newTicketRepo() *TicketRepositoryLocal {
	return &TicketRepositoryLocal{}
}

func newTicketOpen() *domain.Ticket {
	title, description := "title", "description"
	client, _ := domain.UserRegister("client", "client@mail.com", "password")
	ticket, _ := domain.TicketOpen(title, description, client)
	return ticket
}
func TestTicketInsert(t *testing.T) {
	ticket := newTicketOpen()

	t.Run("a valid ticket", func(t *testing.T) {
		repo := newTicketRepo()
		got, _ := repo.Insert(ticket)

		if got == nil {
			t.Error("expected a ticket, but got a nil")
		}
	})

	t.Run("a duplicated ticket", func(t *testing.T) {
		repo := newTicketRepo()
		repo.Insert(ticket)
		_, err := repo.Insert(ticket)
		if err == nil {
			t.Error("expected an error, but got a nil")
		}
	})
}

func TestGetById(t *testing.T) {
	ticket := newTicketOpen()

	t.Run("get an existent ticket", func(t *testing.T) {
		repo := newTicketRepo()
		repo.Insert(ticket)
		got, _ := repo.GetById(ticket.ID)
		if got == nil {
			t.Error("expected a ticket, but got a nil")
		}
	})

	t.Run("get a no-existent ticket", func(t *testing.T) {
		repo := newTicketRepo()
		_, err := repo.GetById(ticket.ID)
		if err == nil {
			t.Error("expected an error, but got a nil")
		}

		if err != nil && err.Type != myerrors.REGISTER_NOT_FOUND {
			t.Errorf("expected %q, but got %q", myerrors.REGISTER_NOT_FOUND, err.Type)
		}
	})
}

func TestTicketUpdate(t *testing.T) {
	operator, _ := domain.UserRegister("operator", "operator@mail.com", "password")
	operator.Profile = domain.PROFILE_OPERATOR

	t.Run("update a valid ticket", func(t *testing.T) {
		repo := newTicketRepo()
		ticket := newTicketOpen()
		repo.Insert(ticket)
		ticket.Get(operator)
		got, _ := repo.Update(ticket)

		if got == nil {
			t.Error("expected a ticket, but got a nil")
		}

		if len(repo.Tickets) != 1 {
			t.Errorf("expected %q tickets, but got %q", 1, len(repo.Tickets))
		}
	})

	t.Run("update a no-existent ticket", func(t *testing.T) {
		repo := newTicketRepo()
		ticket := newTicketOpen()
		ticket.Get(operator)
		_, err := repo.Update(ticket)

		if err == nil {
			t.Error("expected an error, but got a nil")
		}

		if len(repo.Tickets) > 0 {
			t.Errorf("expected %q tickets, but got %q", 1, len(repo.Tickets))
		}
	})
}

func TestTicketGetAll(t *testing.T) {
	operator, _ := domain.UserRegister("operator", "operator@mail.com", "password")
	operator.Profile = domain.PROFILE_OPERATOR
	ticketOpen := newTicketOpen()
	ticketInProgress := newTicketOpen()
	ticketInProgress.Get(operator)

	t.Run("get all tickets", func(t *testing.T) {
		repo := newTicketRepo()
		repo.Insert(ticketOpen)
		repo.Insert(ticketInProgress)
		tickets, _ := repo.GetAll()
		if len(*tickets) != 2 {
			t.Errorf("expected %q, but got %q", 2, len(*tickets))
		}
	})
}

func TestTicketGetAllOpen(t *testing.T) {
	operator, _ := domain.UserRegister("operator", "operator@mail.com", "password")
	operator.Profile = domain.PROFILE_OPERATOR
	ticketOpen := newTicketOpen()
	ticketInProgress := newTicketOpen()
	ticketInProgress.Get(operator)

	t.Run("get all open tickets", func(t *testing.T) {
		repo := newTicketRepo()
		repo.Insert(ticketOpen)
		repo.Insert(ticketInProgress)
		tickets, _ := repo.GetAllOpen()
		if len(*tickets) != 1 {
			t.Errorf("expected %q, but got %q", 1, len(*tickets))
		}
	})
}

func TestTicketGetAllByOperator(t *testing.T) {
	operator, _ := domain.UserRegister("operator", "operator@mail.com", "password")
	operator.Profile = domain.PROFILE_OPERATOR
	ticketOpen := newTicketOpen()
	ticketInProgress := newTicketOpen()
	ticketInProgress.Get(operator)

	t.Run("get all tickets by operator", func(t *testing.T) {
		repo := newTicketRepo()
		repo.Insert(ticketOpen)
		repo.Insert(ticketInProgress)
		tickets, _ := repo.GetAllByOperator(operator)
		if len(*tickets) != 1 {
			t.Errorf("expected %q, but got %q", 1, len(*tickets))
		}
	})
}

func TestTicketGetAllByClient(t *testing.T) {
	ticket1 := newTicketOpen()
	client, _ := domain.UserRegister("client2", "client2@mail.com", "password")
	ticket2, _ := domain.TicketOpen("title", "description", client)

	t.Run("get all tickets by client", func(t *testing.T) {
		repo := newTicketRepo()
		repo.Insert(ticket1)
		repo.Insert(ticket2)
		tickets, _ := repo.GetAllByClient(client)
		if len(*tickets) != 1 {
			t.Errorf("expected %d, but got %d", 1, len(*tickets))
		}
	})
}

func TestTicketDelete(t *testing.T) {
	ticket := newTicketOpen()

	t.Run("delete an existent ticket", func(t *testing.T) {
		repo := newTicketRepo()
		repo.Insert(ticket)
		err := repo.Delete(ticket.ID)
		if err != nil {
			t.Error("expected a nil, but got an error")
		}
	})

	t.Run("delete a no-existent ticket", func(t *testing.T) {
		repo := newTicketRepo()
		err := repo.Delete(ticket.ID)
		if err == nil {
			t.Error("expected an error, but got a nil")
		}
	})
}
