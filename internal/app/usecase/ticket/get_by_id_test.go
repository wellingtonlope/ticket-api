package ticket

import (
	"testing"
	"time"

	"github.com/wellingtonlope/ticket-api/internal/app/security"

	"github.com/stretchr/testify/assert"
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/domain"
	"github.com/wellingtonlope/ticket-api/internal/infra/memory"
)

func TestGetByID(t *testing.T) {
	t.Run("Should get a ticket by ID", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewGetByID(repo)

		client, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
		client, _ = repoUser.Insert(*client)
		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
		operator.Profile = domain.ProfileOperator
		operator, _ = repoUser.Insert(*operator)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *client)
		updatedTicket, _ := ticket.Get(*operator, time.Now())
		insertedTicket, _ := repo.Insert(updatedTicket)

		input := GetByIDInput{
			TicketID:   insertedTicket.ID,
			LoggedUser: security.NewUser(*client),
		}

		output, err := uc.Handle(input)
		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, insertedTicket.ID, output.ID)
		assert.Equal(t, insertedTicket.Title, output.Title)
		assert.Equal(t, insertedTicket.Description, output.Description)
		assert.Equal(t, insertedTicket.CreatedAt, output.CreatedAt)
		assert.Equal(t, insertedTicket.UpdatedAt, output.UpdatedAt)
		assert.Equal(t, insertedTicket.Client.ID, output.Client.ID)
		assert.Equal(t, insertedTicket.Client.Name, output.Client.Name)
		assert.Equal(t, insertedTicket.Client.Email.String(), output.Client.Email)
		assert.Equal(t, string(insertedTicket.Status), output.Status)
		assert.Equal(t, insertedTicket.Operator.ID, output.Operator.ID)
		assert.Equal(t, insertedTicket.Operator.Name, output.Operator.Name)
		assert.Equal(t, insertedTicket.Operator.Email.String(), output.Operator.Email)
	})

	t.Run("Shouldn't get a ticket by another user", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewGetByID(repo)

		client, _ := domain.UserRegister("client1", "client1@mail.com", "password", time.Now())
		client, _ = repoUser.Insert(*client)
		clientOther, _ := domain.UserRegister("client2", "client2@mail.com", "password", time.Now())
		clientOther, _ = repoUser.Insert(*clientOther)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *clientOther)
		insertedTicket, _ := repo.Insert(ticket)

		input := GetByIDInput{
			TicketID:   insertedTicket.ID,
			LoggedUser: security.NewUser(*client),
		}

		output, err := uc.Handle(input)
		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, repository.ErrTicketNotFound, err)
	})

	t.Run("Should get a ticket by another user when operator user", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewGetByID(repo)

		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
		operator.Profile = domain.ProfileOperator
		operator, _ = repoUser.Insert(*operator)
		client, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
		client, _ = repoUser.Insert(*client)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *client)
		insertedTicket, _ := repo.Insert(ticket)

		input := GetByIDInput{
			TicketID:   insertedTicket.ID,
			LoggedUser: security.NewUser(*operator),
		}

		output, err := uc.Handle(input)
		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, insertedTicket.ID, output.ID)
	})

	t.Run("Shouldn't get a ticket by ID when invalid input", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewGetByID(repo)

		client, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
		client, _ = repoUser.Insert(*client)

		input := GetByIDInput{
			TicketID:   "invalid_id",
			LoggedUser: security.NewUser(*client),
		}

		output, err := uc.Handle(input)
		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, repository.ErrTicketNotFound, err)
	})
}
