package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wellingtonlope/ticket-api/internal/app"
	"github.com/wellingtonlope/ticket-api/internal/domain"
	"github.com/wellingtonlope/ticket-api/internal/infra/memory"
)

func TestGetAllOperators(t *testing.T) {
	operator1, _ := domain.UserRegister("operator1", "operator1@mail.com", "password", time.Now())
	operator1.Profile = domain.PROFILE_OPERATOR
	operator2, _ := domain.UserRegister("operator2", "operator2@mail.com", "password", time.Now())
	operator2.Profile = domain.PROFILE_OPERATOR
	client, _ := domain.UserRegister("client", "client@mail.com", "password", time.Now())

	t.Run("should return all operators", func(t *testing.T) {
		repo := &memory.UserRepository{}
		uc := NewGetAllOperators(repo)

		operator1, _ = repo.Insert(*operator1)
		repo.Insert(*operator2)
		repo.Insert(*client)

		output, err := uc.Handle(GetAllOperatorsInput{LoggedUser: *operator1})

		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Len(t, *output, 2)
	})

	t.Run("shouldn't get all operators when user don't have permission", func(t *testing.T) {
		repo := &memory.UserRepository{}
		uc := NewGetAllOperators(repo)

		repo.Insert(*operator1)
		repo.Insert(*operator2)
		client, _ := repo.Insert(*client)

		output, err := uc.Handle(GetAllOperatorsInput{LoggedUser: *client})

		assert.Nil(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, app.ErrForbidden, err)
	})
}
