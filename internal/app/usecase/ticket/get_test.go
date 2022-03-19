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

func TestGet(t *testing.T) {
	t.Run("should get a ticket", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewGet(repo, repoUser)

		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
		operator.Profile = domain.PROFILE_OPERATOR
		operator, _ = repoUser.Insert(*operator)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *operator)
		ticket, _ = repo.Insert(*ticket)
		expectedUpdatedAt := time.Now()

		input := GetInput{TicketID: ticket.ID, UpdatedAt: expectedUpdatedAt, LoggedUser: security.NewUser(*operator)}
		output, err := uc.Handle(input)

		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, ticket.ID, output.ID)
		assert.Equal(t, string(domain.STATUS_IN_PROGRESS), output.Status)
		assert.Equal(t, expectedUpdatedAt, *output.UpdatedAt)
		assert.Equal(t, operator.ID, output.Operator.ID)
		assert.Equal(t, operator.Name, output.Operator.Name)
		assert.Equal(t, operator.Email.String(), output.Operator.Email)

		ticketRepo, _ := repo.GetByID(ticket.ID)
		assert.Equal(t, domain.STATUS_IN_PROGRESS, ticketRepo.Status)
		assert.Equal(t, expectedUpdatedAt, *ticketRepo.UpdatedAt)
		assert.Equal(t, operator.ID, ticketRepo.Operator.ID)
		assert.Equal(t, operator.Name, ticketRepo.Operator.Name)
		assert.Equal(t, operator.Email.String(), ticketRepo.Operator.Email.String())

	})

	t.Run("should return an error when user is a client", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewGet(repo, repoUser)

		client, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
		client, _ = repoUser.Insert(*client)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *client)
		ticket, _ = repo.Insert(*ticket)

		input := GetInput{TicketID: ticket.ID, LoggedUser: security.NewUser(*client)}
		output, err := uc.Handle(input)

		assert.Nil(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, security.ErrForbidden, err)
	})

	t.Run("should return an error when ticket not exists", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewGet(repo, repoUser)

		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
		operator.Profile = domain.PROFILE_OPERATOR
		operator, _ = repoUser.Insert(*operator)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *operator)
		ticket, _ = repo.Insert(*ticket)
		expectedUpdatedAt := time.Now()

		input := GetInput{TicketID: "invalid-id", UpdatedAt: expectedUpdatedAt, LoggedUser: security.NewUser(*operator)}
		output, err := uc.Handle(input)

		assert.Nil(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, repository.ErrTicketNotFound, err)
	})
}
