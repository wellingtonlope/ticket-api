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

func TestAssignToOperator(t *testing.T) {
	t.Run("should assign ticket to an operator", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewAssignToOperator(repo, repoUser)

		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
		operator.Profile = domain.PROFILE_OPERATOR
		operator, _ = repoUser.Insert(*operator)
		operatorOther, _ := domain.UserRegister("operatorOther", "operatorOther@mail.com", "password", time.Now())
		operatorOther.Profile = domain.PROFILE_OPERATOR
		operatorOther, _ = repoUser.Insert(*operatorOther)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *operatorOther)
		ticket, _ = repo.Insert(*ticket)
		expectedUpdatedAt := time.Now()

		input := AssignToOperatorInput{TicketID: ticket.ID, OperatorID: operatorOther.ID, UpdatedAt: expectedUpdatedAt, LoggedUser: security.NewUser(*operator)}
		output, err := uc.Handle(input)

		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, ticket.ID, output.ID)
		assert.Equal(t, string(domain.STATUS_IN_PROGRESS), output.Status)
		assert.Equal(t, expectedUpdatedAt, *output.UpdatedAt)
		assert.Equal(t, operatorOther.ID, output.Operator.ID)
		assert.Equal(t, operatorOther.Name, output.Operator.Name)
		assert.Equal(t, operatorOther.Email.String(), output.Operator.Email)

		ticketRepo, _ := repo.GetByID(ticket.ID)
		assert.Equal(t, domain.STATUS_IN_PROGRESS, ticketRepo.Status)
		assert.Equal(t, expectedUpdatedAt, *ticketRepo.UpdatedAt)
		assert.Equal(t, operatorOther.ID, ticketRepo.Operator.ID)
		assert.Equal(t, operatorOther.Name, ticketRepo.Operator.Name)
		assert.Equal(t, operatorOther.Email.String(), ticketRepo.Operator.Email.String())
	})

	t.Run("should return an error when user assigned was a client", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewAssignToOperator(repo, repoUser)

		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
		operator.Profile = domain.PROFILE_OPERATOR
		operator, _ = repoUser.Insert(*operator)
		client, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
		client, _ = repoUser.Insert(*client)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *client)
		ticket, _ = repo.Insert(*ticket)
		expectedUpdatedAt := time.Now()

		input := AssignToOperatorInput{TicketID: ticket.ID, OperatorID: client.ID, UpdatedAt: expectedUpdatedAt, LoggedUser: security.NewUser(*operator)}
		output, err := uc.Handle(input)

		assert.Nil(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, domain.ErrTicketNoOperator, err)
	})

	t.Run("should return an error when logged user is client", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewAssignToOperator(repo, repoUser)

		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
		operator.Profile = domain.PROFILE_OPERATOR
		operator, _ = repoUser.Insert(*operator)
		client, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
		client, _ = repoUser.Insert(*client)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *client)
		ticket, _ = repo.Insert(*ticket)
		expectedUpdatedAt := time.Now()

		input := AssignToOperatorInput{TicketID: ticket.ID, OperatorID: operator.ID, UpdatedAt: expectedUpdatedAt, LoggedUser: security.NewUser(*client)}
		output, err := uc.Handle(input)

		assert.Nil(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, security.ErrForbidden, err)
	})

	t.Run("should return an error when operator not exists", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewAssignToOperator(repo, repoUser)

		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
		operator.Profile = domain.PROFILE_OPERATOR
		operator, _ = repoUser.Insert(*operator)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *operator)
		ticket, _ = repo.Insert(*ticket)
		expectedUpdatedAt := time.Now()

		input := AssignToOperatorInput{TicketID: ticket.ID, OperatorID: "invalid-id", UpdatedAt: expectedUpdatedAt, LoggedUser: security.NewUser(*operator)}
		output, err := uc.Handle(input)

		assert.Nil(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, repository.ErrUserNotFound, err)
	})

	t.Run("should return an error when ticket not exists", func(t *testing.T) {
		repo := &memory.TicketRepository{}
		repoUser := &memory.UserRepository{}
		uc := NewAssignToOperator(repo, repoUser)

		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
		operator.Profile = domain.PROFILE_OPERATOR
		operator, _ = repoUser.Insert(*operator)
		ticket, _ := domain.OpenTicket("title", "description", time.Now(), *operator)
		ticket, _ = repo.Insert(*ticket)
		expectedUpdatedAt := time.Now()

		input := AssignToOperatorInput{TicketID: "invalid-id", OperatorID: operator.ID, UpdatedAt: expectedUpdatedAt, LoggedUser: security.NewUser(*operator)}
		output, err := uc.Handle(input)

		assert.Nil(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, repository.ErrTicketNotFound, err)
	})
}
