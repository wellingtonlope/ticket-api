package ticket

import (
	"testing"
	"time"

	"github.com/wellingtonlope/ticket-api/internal/app/security"

	"github.com/stretchr/testify/assert"
	"github.com/wellingtonlope/ticket-api/internal/domain"
	"github.com/wellingtonlope/ticket-api/internal/infra/memory"
)

func TestOpen(t *testing.T) {
	client, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())

	t.Run("should open ticket", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewOpen(repo, repoUser)

		client, _ := repoUser.Insert(*client)

		input := OpenInput{
			Title:       "title",
			Description: "description",
			CreatedAt:   time.Now(),
			LoggedUser:  security.NewUser(*client),
		}

		output, err := uc.Handle(input)
		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.NotEmpty(t, output.ID)
		assert.Equal(t, input.Title, output.Title)
		assert.Equal(t, input.Description, output.Description)
		assert.Equal(t, input.CreatedAt, *output.CreatedAt)
		assert.Equal(t, string(domain.StatusOpen), output.Status)
		assert.Equal(t, input.LoggedUser.ID, output.Client.ID)
		assert.Equal(t, input.LoggedUser.Name, output.Client.Name)

		ticket, err := repo.GetByID(output.ID)
		assert.Nil(t, err)
		assert.NotNil(t, ticket)
		assert.Equal(t, output.ID, ticket.ID)
		assert.Equal(t, output.Title, ticket.Title)
		assert.Equal(t, output.Description, ticket.Description)
		assert.Equal(t, output.CreatedAt, ticket.CreatedAt)
		assert.Equal(t, domain.StatusOpen, ticket.Status)
		assert.Equal(t, output.Client.ID, ticket.Client.ID)
		assert.Equal(t, output.Client.Name, ticket.Client.Name)
		assert.Equal(t, output.Client.Email, ticket.Client.Email.String())
	})

	t.Run("shouldn't open ticket when invalid input", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewOpen(repo, repoUser)

		client, _ := repoUser.Insert(*client)

		input := OpenInput{
			Title:       "",
			Description: "",
			CreatedAt:   time.Now(),
			LoggedUser:  security.NewUser(*client),
		}

		output, err := uc.Handle(input)
		assert.Nil(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, domain.ErrTicketTitleIsInvalid, err)
	})
}
