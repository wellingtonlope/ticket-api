package ticket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/domain"
	"github.com/wellingtonlope/ticket-api/internal/infra/memory"
)

func TestGetAllOpen(t *testing.T) {
	t.Run("should return all open tickets", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewGetAllOpen(repo)

		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
		operator.Profile = domain.ProfileOperator
		operator, _ = repoUser.Insert(*operator)
		ticketInProgress, _ := domain.OpenTicket("title1", "description1", time.Now(), *operator)
		insertedTicket, _ := repo.Insert(ticketInProgress)
		_ = insertedTicket.Get(*operator, time.Now())
		_, _ = repo.Update(*insertedTicket)
		ticketOpen, _ := domain.OpenTicket("title2", "description2", time.Now(), *operator)
		_, _ = repo.Insert(ticketOpen)

		input := GetAllOpenInput{LoggedUser: security.NewUser(*operator)}
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
		operator.Profile = domain.ProfileOperator
		operator, _ = repoUser.Insert(*operator)
		client, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
		client, _ = repoUser.Insert(*client)
		ticketInProgress, _ := domain.OpenTicket("title", "description", time.Now(), *client)
		insertedTicket, _ := repo.Insert(ticketInProgress)
		_ = insertedTicket.Get(*operator, time.Now())
		_, _ = repo.Update(*insertedTicket)
		ticketOpen, _ := domain.OpenTicket("title2", "description2", time.Now(), *client)
		_, _ = repo.Insert(ticketOpen)

		input := GetAllOpenInput{LoggedUser: security.NewUser(*client)}
		output, err := uc.Handle(input)

		assert.Nil(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, security.ErrForbidden, err)
	})
}
