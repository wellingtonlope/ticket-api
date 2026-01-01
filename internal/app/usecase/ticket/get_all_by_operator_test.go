package ticket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/domain"
	"github.com/wellingtonlope/ticket-api/internal/infra/memory"
)

func TestGetAllByOperator(t *testing.T) {
	t.Run("should return all tickets by operator", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewGetAllByOperator(repo)

		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
		operator.Profile = domain.ProfileOperator
		operator, _ = repoUser.Insert(*operator)
		operatorOther, _ := domain.UserRegister("operatorOther", "operatorOther@mail.com", "password", time.Now())
		operatorOther, _ = repoUser.Insert(*operatorOther)
		ticketOther, _ := domain.OpenTicket("title", "description", time.Now(), *operatorOther)
		updatedTicketOther, _ := ticketOther.Get(*operatorOther, time.Now())
		_, _ = repo.Insert(updatedTicketOther)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *operator)
		updatedTicket, _ := ticket.Get(*operator, time.Now())
		_, _ = repo.Insert(updatedTicket)

		input := GetAllByOperatorInput{OperatorID: operator.ID, LoggedUser: security.NewUser(*operator)}
		output, err := uc.Handle(input)

		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Len(t, *output, 1)
	})

	t.Run("should return error logged user is a client", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewGetAllByOperator(repo)

		client, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
		client, _ = repoUser.Insert(*client)
		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
		operator.Profile = domain.ProfileOperator
		operator, _ = repoUser.Insert(*operator)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *operator)
		updatedTicket, _ := ticket.Get(*operator, time.Now())
		_, _ = repo.Insert(updatedTicket)

		input := GetAllByOperatorInput{OperatorID: operator.ID, LoggedUser: security.NewUser(*client)}
		output, err := uc.Handle(input)

		assert.Nil(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, security.ErrForbidden, err)
	})
}
