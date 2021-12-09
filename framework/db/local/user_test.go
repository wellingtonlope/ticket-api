package local

import (
	"testing"

	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
)

func newUserRepo() *UserRepositoryLocal {
	return &UserRepositoryLocal{}
}
func TestUserInsert(t *testing.T) {
	user, _ := domain.UserRegister("user", "user@mail.com", "password")

	t.Run("a valid user", func(t *testing.T) {
		repo := newUserRepo()
		got, _ := repo.Insert(user)
		if got == nil {
			t.Error("expected an user, but got a nil")
		}
	})

	t.Run("a duplicated user", func(t *testing.T) {
		repo := newUserRepo()
		repo.Insert(user)
		_, err := repo.Insert(user)
		if err == nil {
			t.Error("expected an error, but got a nil")
		}
	})
}

func TestUserGetById(t *testing.T) {
	user, _ := domain.UserRegister("user", "user@mail.com", "password")

	t.Run("get an existent user", func(t *testing.T) {
		repo := newUserRepo()
		repo.Insert(user)
		got, _ := repo.GetById(user.ID)
		if got == nil {
			t.Error("expected an user, but got a nil")
		}
	})

	t.Run("get a no-existent user", func(t *testing.T) {
		repo := newUserRepo()
		_, err := repo.GetById(user.ID)
		if err == nil {
			t.Error("expected an error, but got a nil")
		}

		if err != nil && err.Type != myerrors.REGISTER_NOT_FOUND {
			t.Errorf("expected %q, but got %q", myerrors.REGISTER_NOT_FOUND, err.Type)
		}
	})
}

func TestUserGetByEmail(t *testing.T) {
	user, _ := domain.UserRegister("user", "user@mail.com", "password")

	t.Run("get an existent user", func(t *testing.T) {
		repo := newUserRepo()
		repo.Insert(user)
		got, _ := repo.GetByEmail(user.Email)
		if got == nil {
			t.Error("expected an user, but got a nil")
		}
	})

	t.Run("get a no-existent user", func(t *testing.T) {
		repo := newUserRepo()
		_, err := repo.GetByEmail(user.Email)
		if err == nil {
			t.Error("expected an error, but got a nil")
		}

		if err != nil && err.Type != myerrors.REGISTER_NOT_FOUND {
			t.Errorf("expected %q, but got %q", myerrors.REGISTER_NOT_FOUND, err.Type)
		}
	})
}

func TestUserGetAllOperator(t *testing.T) {
	client, _ := domain.UserRegister("client", "client@mail.com", "password")
	operator, _ := domain.UserRegister("operator", "operator@mail.com", "password")
	operator.Profile = domain.PROFILE_OPERATOR
	repo := newUserRepo()
	repo.Insert(client)
	repo.Insert(operator)

	t.Run("get all operators", func(t *testing.T) {
		users, _ := repo.GetAllOperator()
		if len(*users) != 1 {
			t.Errorf("expected %q, but got %q", 1, len(*users))
		}
	})
}

func TestUserDelete(t *testing.T) {
	user, _ := domain.UserRegister("user", "user@mail.com", "password")

	t.Run("delete an existent user", func(t *testing.T) {
		repo := newUserRepo()
		repo.Insert(user)
		err := repo.Delete(user.ID)
		if err != nil {
			t.Error("expected a nil, but got an error")
		}
	})

	t.Run("delete a no-existent user", func(t *testing.T) {
		repo := newUserRepo()
		err := repo.Delete(user.ID)
		if err == nil {
			t.Error("expected an error, but got a nil")
		}
	})
}
