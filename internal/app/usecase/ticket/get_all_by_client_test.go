package ticket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/domain"
	"github.com/wellingtonlope/ticket-api/internal/infra/memory"
)

func TestGetAllByClient(t *testing.T) {
	t.Run("should return all tickets by client", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewGetAllByClient(repo)

		client, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
		client, _ = repoUser.Insert(*client)
		clientOther, _ := domain.UserRegister("clientOther", "clientOther@mail.com", "password", time.Now())
		clientOther, _ = repoUser.Insert(*clientOther)
		_, _ = domain.OpenTicket("title", "description", time.Now(), *clientOther)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *client)
		ticket, _ = repo.Insert(*ticket)

		input := GetAllByClientInput{ClientID: client.ID, LoggedUser: security.NewUser(*client)}
		output, err := uc.Handle(input)

		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Len(t, *output, 1)
	})

	t.Run("should return all tickets by client when logged user is an operator", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewGetAllByClient(repo)

		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
		operator.Profile = domain.PROFILE_OPERATOR
		operator, _ = repoUser.Insert(*operator)
		clientOther, _ := domain.UserRegister("clientOther", "clientOther@mail.com", "password", time.Now())
		clientOther, _ = repoUser.Insert(*clientOther)
		ticketOther, _ := domain.OpenTicket("title", "description", time.Now(), *clientOther)
		ticketOther, _ = repo.Insert(*ticketOther)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *operator)
		ticket, _ = repo.Insert(*ticket)

		input := GetAllByClientInput{ClientID: clientOther.ID, LoggedUser: security.NewUser(*operator)}
		output, err := uc.Handle(input)

		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Len(t, *output, 1)
	})

	t.Run("should return error when clientID is different from logged user", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewGetAllByClient(repo)

		client, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
		client, _ = repoUser.Insert(*client)
		clientOther, _ := domain.UserRegister("clientOther", "clientOther@mail.com", "password", time.Now())
		clientOther, _ = repoUser.Insert(*clientOther)
		_, _ = domain.OpenTicket("title", "description", time.Now(), *clientOther)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *client)
		ticket, _ = repo.Insert(*ticket)

		input := GetAllByClientInput{ClientID: clientOther.ID, LoggedUser: security.NewUser(*client)}
		output, err := uc.Handle(input)

		assert.Nil(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, security.ErrForbidden, err)
	})
}
