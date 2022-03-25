package memory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

func TestUserRepository_GetByID(t *testing.T) {
	userPassword := "password"
	userFixture, _ := domain.UserRegister("name", "email@mail.com", userPassword, time.Now())

	t.Run("should get a user by id", func(t *testing.T) {
		repo := &UserRepository{}
		user, _ := repo.Insert(*userFixture)

		got, err := repo.GetByID(user.ID)

		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, user.ID, got.ID)
		assert.Equal(t, user.Name, got.Name)
		assert.Equal(t, user.Email.String(), got.Email.String())
		assert.True(t, got.Password.IsCorrectPassword(userPassword))
		assert.Equal(t, user.Profile, got.Profile)
		assert.Equal(t, user.CreatedAt, got.CreatedAt)
		assert.Equal(t, user.UpdatedAt, got.UpdatedAt)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		repo := &UserRepository{}

		got, err := repo.GetByID("ID_NOT_FOUND")

		assert.NotNil(t, err)
		assert.Nil(t, got)
		assert.Equal(t, repository.ErrUserNotFound, err)
	})
}

func TestUserRepository_Insert(t *testing.T) {
	userPassword := "password"
	userFixture, _ := domain.UserRegister("name", "email@mail.com", userPassword, time.Now())

	t.Run("should insert a user", func(t *testing.T) {
		repo := &UserRepository{}

		got, err := repo.Insert(*userFixture)

		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.NotEmpty(t, got.ID)
		assert.Equal(t, userFixture.Name, got.Name)
		assert.Equal(t, userFixture.Email.String(), got.Email.String())
		assert.True(t, got.Password.IsCorrectPassword(userPassword))
		assert.Equal(t, userFixture.Profile, got.Profile)
		assert.Equal(t, *userFixture.CreatedAt, *got.CreatedAt)

		gotRepo, err := repo.GetByID(got.ID)

		assert.Nil(t, err)
		assert.NotNil(t, gotRepo)
		assert.Equal(t, got.ID, gotRepo.ID)
		assert.Equal(t, got.Name, gotRepo.Name)
		assert.Equal(t, got.Email.String(), gotRepo.Email.String())
		assert.True(t, gotRepo.Password.IsCorrectPassword(userPassword))
		assert.Equal(t, got.Profile, gotRepo.Profile)
		assert.Equal(t, *got.CreatedAt, *gotRepo.CreatedAt)
	})
}

func TestUserRepository_GetByEmail(t *testing.T) {
	userPassword := "password"
	userFixture, _ := domain.UserRegister("name", "email@mail.com", userPassword, time.Now())

	t.Run("should get a user by email", func(t *testing.T) {
		repo := &UserRepository{}
		user, _ := repo.Insert(*userFixture)

		got, err := repo.GetByEmail(user.Email.String())

		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, user.ID, got.ID)
		assert.Equal(t, user.Name, got.Name)
		assert.Equal(t, user.Email.String(), got.Email.String())
		assert.True(t, got.Password.IsCorrectPassword(userPassword))
		assert.Equal(t, user.Profile, got.Profile)
		assert.Equal(t, user.CreatedAt, got.CreatedAt)
		assert.Equal(t, user.UpdatedAt, got.UpdatedAt)
	})

	t.Run("should return error not found", func(t *testing.T) {
		repo := &UserRepository{}

		got, err := repo.GetByEmail("not_found@mail.com")

		assert.NotNil(t, err)
		assert.Nil(t, got)
		assert.Equal(t, repository.ErrUserNotFound, err)
	})
}

func TestUserRepository_GetAllOperator(t *testing.T) {
	userOperatorFixture, _ := domain.UserRegister("name", "operator@mail.com", "password", time.Now())
	userOperatorFixture.Profile = domain.ProfileOperator
	userClientFixture, _ := domain.UserRegister("name", "client@mail.com", "password", time.Now())

	t.Run("should get all operator users", func(t *testing.T) {
		repo := &UserRepository{}
		_, _ = repo.Insert(*userOperatorFixture)
		_, _ = repo.Insert(*userOperatorFixture)
		_, _ = repo.Insert(*userClientFixture)

		got, err := repo.GetAllOperator()

		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Len(t, *got, 2)
	})
}
