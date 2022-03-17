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

func TestClose(t *testing.T) {
	t.Run("should close a ticket", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewClose(repo)

		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
		operator.Profile = domain.PROFILE_OPERATOR
		operator, _ = repoUser.Insert(*operator)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *operator)
		ticket.Get(*operator, time.Now())
		ticket, _ = repo.Insert(*ticket)

		input := CloseInput{TicketID: ticket.ID, Solution: "solution", UpdatedAt: time.Now(), LoggedUser: *operator}
		output, err := uc.Handle(input)

		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, ticket.ID, output.ID)
		assert.Equal(t, string(domain.STATUS_CLOSE), output.Status)
		assert.Equal(t, input.UpdatedAt, *output.UpdatedAt)
		assert.Equal(t, input.Solution, output.Solution)

		ticketRepo, _ := repo.GetByID(ticket.ID)
		assert.Equal(t, domain.STATUS_CLOSE, ticketRepo.Status)
		assert.Equal(t, input.UpdatedAt, *ticketRepo.UpdatedAt)
		assert.Equal(t, input.Solution, output.Solution)
	})

	t.Run("should return an error when user is a client", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewClose(repo)

		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
		operator.Profile = domain.PROFILE_OPERATOR
		operator, _ = repoUser.Insert(*operator)
		client, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
		client, _ = repoUser.Insert(*client)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *client)
		ticket.Get(*operator, time.Now())
		ticket, _ = repo.Insert(*ticket)

		input := CloseInput{TicketID: ticket.ID, Solution: "solution", UpdatedAt: time.Now(), LoggedUser: *client}
		output, err := uc.Handle(input)

		assert.Nil(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, security.ErrForbidden, err)
	})

	t.Run("should return an error when ticket not exists", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewClose(repo)

		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
		operator.Profile = domain.PROFILE_OPERATOR
		operator, _ = repoUser.Insert(*operator)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *operator)
		ticket.Get(*operator, time.Now())
		ticket, _ = repo.Insert(*ticket)

		input := CloseInput{TicketID: "invalid-id", Solution: "solution", UpdatedAt: time.Now(), LoggedUser: *operator}
		output, err := uc.Handle(input)

		assert.Nil(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, repository.ErrTicketNotFound, err)
	})
}
