package user

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
	"testing"
)

func TestGetById(t *testing.T) {
	uc := newUserUC()
	user, _ := domain.UserRegister("client", "client@mail.com", "password")
	_, _ = uc.UserRepository.Insert(user)

	t.Run("get an existent user", func(t *testing.T) {
		got, _ := uc.GetById(user.ID)

		if got == nil {
			t.Error("expected an user, but got a nil")
		}
	})

	t.Run("get a no-existent user", func(t *testing.T) {
		_, myerr := uc.GetById("123")

		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.REGISTER_NOT_FOUND {
			t.Errorf("expected %q, but got %q", myerrors.REGISTER_NOT_FOUND, myerr.Type)
		}
	})
}
