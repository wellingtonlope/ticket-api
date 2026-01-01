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
		operator.Profile = domain.ProfileOperator
		operator, _ = repoUser.Insert(*operator)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *operator)
		updatedTicket, _ := ticket.Get(*operator, time.Now())
		insertedTicket, _ := repo.Insert(updatedTicket)

		input := CloseInput{TicketID: insertedTicket.ID, Solution: "solution", UpdatedAt: time.Now(), LoggedUser: security.NewUser(*operator)}
		output, err := uc.Handle(input)

		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, insertedTicket.ID, output.ID)
		assert.Equal(t, string(domain.StatusClose), output.Status)
		assert.Equal(t, input.UpdatedAt, output.UpdatedAt)
		assert.Equal(t, input.Solution, output.Solution)

		ticketRepo, _ := repo.GetByID(insertedTicket.ID)
		assert.Equal(t, domain.StatusClose, ticketRepo.Status)
		assert.Equal(t, input.UpdatedAt, ticketRepo.UpdatedAt)
		assert.Equal(t, input.Solution, output.Solution)
	})

	t.Run("should return an error when user is a client", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewClose(repo)

		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
		operator.Profile = domain.ProfileOperator
		operator, _ = repoUser.Insert(*operator)
		client, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
		client, _ = repoUser.Insert(*client)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *client)
		updatedTicket, _ := ticket.Get(*operator, time.Now())
		insertedTicket, _ := repo.Insert(updatedTicket)

		input := CloseInput{TicketID: insertedTicket.ID, Solution: "solution", UpdatedAt: time.Now(), LoggedUser: security.NewUser(*client)}
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
		operator.Profile = domain.ProfileOperator
		operator, _ = repoUser.Insert(*operator)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *operator)
		updatedTicket, _ := ticket.Get(*operator, time.Now())
		_, _ = repo.Insert(updatedTicket)

		input := CloseInput{TicketID: "invalid-id", Solution: "solution", UpdatedAt: time.Now(), LoggedUser: security.NewUser(*operator)}
		output, err := uc.Handle(input)

		assert.Nil(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, repository.ErrTicketNotFound, err)
	})
}
