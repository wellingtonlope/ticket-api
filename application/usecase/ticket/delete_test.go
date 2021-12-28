package ticket

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
	"testing"
)

func TestDelete(t *testing.T) {
	uc := newTicketUC()
	operator, _ := domain.UserRegister("operator", "operator@mail.com", "password")
	operator.Profile = domain.PROFILE_OPERATOR
	_, _ = uc.UserUseCase.UserRepository.Insert(operator)
	operatorResponse, _ := uc.UserUseCase.Login("operator@mail.com", "password")
	clientResponse, _ := uc.UserUseCase.Register("client", "client@mail.com", "password")

	t.Run("delete a valid ticket with client", func(t *testing.T) {
		ticket, _ := uc.Open("title", "description", clientResponse.Token)

		myerr := uc.Delete(ticket.ID, clientResponse.Token)

		if myerr != nil {
			t.Error("expected no error, but got an error")
		}
	})

	t.Run("delete a valid ticket with operator", func(t *testing.T) {
		ticket, _ := uc.Open("title", "description", operatorResponse.Token)

		myerr := uc.Delete(ticket.ID, operatorResponse.Token)

		if myerr != nil {
			t.Error("expected no error, but got an error")
		}
	})

	t.Run("delete a valid ticket with another user", func(t *testing.T) {
		ticket, _ := uc.Open("title", "description", clientResponse.Token)

		myerr := uc.Delete(ticket.ID, operatorResponse.Token)

		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.FORBIDDEN {
			t.Errorf("expected %q, but go %q", myerrors.FORBIDDEN, myerr.Type)
		}
	})

	t.Run("deleting a ticket picked up", func(t *testing.T) {
		ticket, _ := uc.Open("title", "description", clientResponse.Token)
		_, _ = uc.Get(ticket.ID, operatorResponse.Token)

		myerr := uc.Delete(ticket.ID, clientResponse.Token)

		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.USECASE {
			t.Errorf("expected %q, but go %q", myerrors.USECASE, myerr.Type)
		}
	})

	t.Run("delete a ticket with an invalid ID", func(t *testing.T) {
		myerr := uc.Delete("123", clientResponse.Token)

		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.REGISTER_NOT_FOUND {
			t.Errorf("expected %q, but go %q", myerrors.REGISTER_NOT_FOUND, myerr.Type)
		}
	})

	t.Run("delete a ticket with an invalid token", func(t *testing.T) {
		ticket, _ := uc.Open("title", "description", clientResponse.Token)

		myerr := uc.Delete(ticket.ID, "123")

		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.UNAUTHORIZED {
			t.Errorf("expected %q, but go %q", myerrors.UNAUTHORIZED, myerr.Type)
		}
	})
}
