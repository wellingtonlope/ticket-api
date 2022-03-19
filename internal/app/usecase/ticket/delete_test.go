package ticket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/domain"
	"github.com/wellingtonlope/ticket-api/internal/infra/memory"
)

func TestDelete(t *testing.T) {

	t.Run("Should delete a ticket", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewDelete(repo)

		client, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
		client, _ = repoUser.Insert(*client)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *client)
		ticket, _ = repo.Insert(*ticket)

		input := DeleteInput{
			TicketID:   ticket.ID,
			LoggedUser: security.NewUser(*client),
		}

		output, err := uc.Handle(input)
		assert.Nil(t, err)
		assert.NotNil(t, output)

		ticketRepo, err := repo.GetByID(ticket.ID)
		assert.Nil(t, ticketRepo)
		assert.NotNil(t, err)
		assert.Equal(t, repository.ErrTicketNotFound, err)
	})

	t.Run("Should return an error when ticket not open", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewDelete(repo)

		client, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
		client, _ = repoUser.Insert(*client)
		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
		operator.Profile = domain.PROFILE_OPERATOR
		operator, _ = repoUser.Insert(*operator)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *client)
		ticket.Get(*operator, time.Now())
		ticket, _ = repo.Insert(*ticket)

		input := DeleteInput{
			TicketID:   ticket.ID,
			LoggedUser: security.NewUser(*client),
		}

		output, err := uc.Handle(input)
		assert.Nil(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, security.ErrForbidden, err)
	})

	t.Run("Should return an error when ticket not open", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewDelete(repo)

		client, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
		client, _ = repoUser.Insert(*client)
		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
		operator.Profile = domain.PROFILE_OPERATOR
		operator, _ = repoUser.Insert(*operator)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *client)
		ticket, _ = repo.Insert(*ticket)

		input := DeleteInput{
			TicketID:   ticket.ID,
			LoggedUser: security.NewUser(*operator),
		}

		output, err := uc.Handle(input)
		assert.Nil(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, security.ErrForbidden, err)
	})
}
