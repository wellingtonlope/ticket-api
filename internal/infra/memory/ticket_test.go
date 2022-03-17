package memory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

func TestTicketRepository_GetByID(t *testing.T) {
	userClientFixture, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
	ticketFixture, _ := domain.OpenTicket("title", "description", time.Now(), *userClientFixture)

	t.Run("should get a ticket by id", func(t *testing.T) {
		repo := &TicketRepository{}
		ticket, _ := repo.Insert(*ticketFixture)

		got, err := repo.GetByID(ticket.ID)

		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, ticket.ID, got.ID)
		assert.Equal(t, ticket.Title, got.Title)
		assert.Equal(t, ticket.Description, got.Description)
		assert.Equal(t, ticket.Status, got.Status)
		assert.Equal(t, ticket.CreatedAt, got.CreatedAt)
		assert.Equal(t, ticket.UpdatedAt, got.UpdatedAt)
		assert.Equal(t, ticket.Client.ID, got.Client.ID)
		assert.Equal(t, ticket.Client.Name, got.Client.Name)
		assert.Equal(t, ticket.Client.Email.String(), got.Client.Email.String())
	})

	t.Run("should return error when ticket not found", func(t *testing.T) {
		repo := &TicketRepository{}

		got, err := repo.GetByID("ID_NOT_FOUND")

		assert.NotNil(t, err)
		assert.Nil(t, got)
		assert.Equal(t, repository.ErrTicketNotFound, err)
	})
}
func TestTicketRepository_Insert(t *testing.T) {
	userClientFixture, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
	ticketFixture, _ := domain.OpenTicket("title", "description", time.Now(), *userClientFixture)

	t.Run("should insert a ticket", func(t *testing.T) {
		repo := &TicketRepository{}

		got, err := repo.Insert(*ticketFixture)

		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.NotEmpty(t, got.ID)
		assert.Equal(t, ticketFixture.Title, got.Title)
		assert.Equal(t, ticketFixture.Description, got.Description)
		assert.Equal(t, ticketFixture.Status, got.Status)
		assert.Equal(t, ticketFixture.CreatedAt, got.CreatedAt)
		assert.Equal(t, ticketFixture.UpdatedAt, got.UpdatedAt)
		assert.Equal(t, ticketFixture.Client.ID, got.Client.ID)
		assert.Equal(t, ticketFixture.Client.Name, got.Client.Name)
		assert.Equal(t, ticketFixture.Client.Email.String(), got.Client.Email.String())

		gotRepo, err := repo.GetByID(got.ID)

		assert.Nil(t, err)
		assert.NotNil(t, gotRepo)
		assert.Equal(t, got, gotRepo)
	})
}

func TestTicketRepository_Update(t *testing.T) {
	userClientFixture, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
	operatorClientFixture, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
	operatorClientFixture.Profile = domain.PROFILE_OPERATOR
	ticketFixture, _ := domain.OpenTicket("title", "description", time.Now(), *userClientFixture)

	t.Run("should update a ticket", func(t *testing.T) {
		repo := &TicketRepository{}
		ticket, _ := repo.Insert(*ticketFixture)

		expectedUpdatedAt := time.Now()
		ticket.Get(*operatorClientFixture, expectedUpdatedAt)

		got, err := repo.Update(*ticket)

		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, ticket, got)

		gotRepo, err := repo.GetByID(got.ID)

		assert.Nil(t, err)
		assert.NotNil(t, gotRepo)
		assert.Equal(t, got, gotRepo)
	})

	t.Run("shouldn't update a ticket", func(t *testing.T) {
		repo := &TicketRepository{}

		got, err := repo.Update(*ticketFixture)

		assert.Nil(t, got)
		assert.NotNil(t, err)
		assert.Equal(t, repository.ErrTicketNotFound, err)
	})
}

func TestTicketRepository_GetAll(t *testing.T) {
	userClientFixture, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
	ticketFixture, _ := domain.OpenTicket("title", "description", time.Now(), *userClientFixture)

	t.Run("should get all tickets", func(t *testing.T) {
		repo := &TicketRepository{}
		_, _ = repo.Insert(*ticketFixture)
		_, _ = repo.Insert(*ticketFixture)

		got, err := repo.GetAll()

		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Len(t, *got, 2)
	})
}

func TestTicketRepository_GetAllOpen(t *testing.T) {
	userClientFixture, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
	operatorFixture, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
	operatorFixture.Profile = domain.PROFILE_OPERATOR
	ticketOpenFixture, _ := domain.OpenTicket("title", "description", time.Now(), *userClientFixture)
	ticketGetFixture, _ := domain.OpenTicket("title", "description", time.Now(), *userClientFixture)
	ticketGetFixture.Get(*operatorFixture, time.Now())

	t.Run("should get all open tickets", func(t *testing.T) {
		repo := &TicketRepository{}
		_, _ = repo.Insert(*ticketOpenFixture)
		_, _ = repo.Insert(*ticketGetFixture)

		got, err := repo.GetAllOpen()

		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Len(t, *got, 1)
	})
}

func TestTicketRepository_GetAllByOperatorID(t *testing.T) {
	userClientFixture, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
	operator1Fixture, _ := domain.UserRegister("operator", "operator@mail.com", "password", time.Now())
	operator1Fixture.Profile = domain.PROFILE_OPERATOR
	operator1Fixture.ID = "operator1"
	operator2Fixture, _ := domain.UserRegister("operator2", "operator2@mail.com", "password", time.Now())
	operator2Fixture.ID = "operator2"
	operator2Fixture.Profile = domain.PROFILE_OPERATOR
	ticket1Fixture, _ := domain.OpenTicket("title", "description", time.Now(), *userClientFixture)
	ticket1Fixture.Get(*operator1Fixture, time.Now())
	ticket2Fixture, _ := domain.OpenTicket("title", "description", time.Now(), *userClientFixture)
	ticket2Fixture.Get(*operator2Fixture, time.Now())

	t.Run("should get all open tickets", func(t *testing.T) {
		repo := &TicketRepository{}

		_, _ = repo.Insert(*ticket1Fixture)
		_, _ = repo.Insert(*ticket2Fixture)

		got, err := repo.GetAllByOperatorID(operator1Fixture.ID)

		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Len(t, *got, 1)
	})
}

func TestTicketRepository_GetAllByClientID(t *testing.T) {
	userClient1Fixture, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
	userClient1Fixture.ID = "client1"
	userClient2Fixture, _ := domain.UserRegister("client2", "client2@mail.com", "password", time.Now())
	userClient2Fixture.ID = "client2"
	ticket1Fixture, _ := domain.OpenTicket("title", "description", time.Now(), *userClient1Fixture)
	ticket2Fixture, _ := domain.OpenTicket("title", "description", time.Now(), *userClient2Fixture)

	t.Run("should get all open tickets", func(t *testing.T) {
		repo := &TicketRepository{}

		_, _ = repo.Insert(*ticket1Fixture)
		_, _ = repo.Insert(*ticket2Fixture)

		got, err := repo.GetAllByClientID(userClient1Fixture.ID)

		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Len(t, *got, 1)
	})
}

func TestTicketRepository_DeleteByID(t *testing.T) {
	userClientFixture, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())
	userClientFixture.ID = "client1"
	ticketFixture, _ := domain.OpenTicket("title", "description", time.Now(), *userClientFixture)

	t.Run("should delete a ticket", func(t *testing.T) {
		repo := &TicketRepository{}
		ticket, _ := repo.Insert(*ticketFixture)

		err := repo.DeleteByID(ticket.ID)

		assert.Nil(t, err)
	})

	t.Run("shouldn't delete a ticket", func(t *testing.T) {
		repo := &TicketRepository{}

		err := repo.DeleteByID("invalid")

		assert.NotNil(t, err)
		assert.Equal(t, repository.ErrTicketNotFound, err)
	})
}
