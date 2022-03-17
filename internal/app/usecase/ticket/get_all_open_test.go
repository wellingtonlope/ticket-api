package ticket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wellingtonlope/ticket-api/internal/app"
	"github.com/wellingtonlope/ticket-api/internal/domain"
	"github.com/wellingtonlope/ticket-api/internal/infra/memory"
)

func TestGetAllOpen(t *testing.T) {
	t.Run("should return all open tickets", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewGetAllOpen(repo)

		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
		operator.Profile = domain.PROFILE_OPERATOR
		operator, _ = repoUser.Insert(*operator)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *operator)
		ticket, _ = repo.Insert(*ticket)
		ticket.Get(*operator, time.Now())
		ticket, _ = repo.Insert(*ticket)

		input := GetAllOpenInput{LoggedUser: *operator}
		output, err := uc.Handle(input)

		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Len(t, *output, 1)
	})

	t.Run("should return an error when user is a client", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewGetAllOpen(repo)

		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
		operator.Profile = domain.PROFILE_OPERATOR
		operator, _ = repoUser.Insert(*operator)
		client, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
		client, _ = repoUser.Insert(*client)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *client)
		ticket, _ = repo.Insert(*ticket)
		ticket.Get(*operator, time.Now())
		ticket, _ = repo.Insert(*ticket)

		input := GetAllOpenInput{LoggedUser: *client}
		output, err := uc.Handle(input)

		assert.Nil(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, app.ErrForbidden, err)
	})
}