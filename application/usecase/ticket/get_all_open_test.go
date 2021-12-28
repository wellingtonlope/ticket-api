package ticket

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
	"testing"
)

func TestGetAllOpen(t *testing.T) {
	t.Run("get all tickets with an operator", func(t *testing.T) {
		uc := newTicketUC()
		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password")
		operator.Profile = domain.PROFILE_OPERATOR
		_, _ = uc.UserUseCase.UserRepository.Insert(operator)
		operatorResponse, _ := uc.UserUseCase.Login("operator@mail.com", "password")
		clientResponse, _ := uc.UserUseCase.Register("client", "client@mail.com", "password")
		_, _ = uc.Open("title", "description", operatorResponse.Token)
		_, _ = uc.Open("title", "description", clientResponse.Token)

		got, _ := uc.GetAllOpen(operatorResponse.Token)
		if got == nil {
			t.Error("expected a ticket slice, but got a nil")
		}
		if len(*got) != 2 {
			t.Errorf("expected %q, but got %q", 2, len(*got))
		}
	})

	t.Run("get all tickets without tickets in database", func(t *testing.T) {
		uc := newTicketUC()
		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password")
		operator.Profile = domain.PROFILE_OPERATOR
		_, _ = uc.UserUseCase.UserRepository.Insert(operator)
		operatorResponse, _ := uc.UserUseCase.Login("operator@mail.com", "password")

		got, _ := uc.GetAllOpen(operatorResponse.Token)
		if got == nil {
			t.Error("expected a ticket slice, but got a nil")
		}
		if len(*got) != 0 {
			t.Errorf("expected %q, but got %q", 0, len(*got))
		}
	})

	t.Run("get all tickets with a client", func(t *testing.T) {
		uc := newTicketUC()
		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password")
		operator.Profile = domain.PROFILE_OPERATOR
		_, _ = uc.UserUseCase.UserRepository.Insert(operator)
		operatorResponse, _ := uc.UserUseCase.Login("operator@mail.com", "password")
		clientResponse, _ := uc.UserUseCase.Register("client", "client@mail.com", "password")
		_, _ = uc.Open("title", "description", operatorResponse.Token)
		_, _ = uc.Open("title", "description", clientResponse.Token)

		_, myerr := uc.GetAllOpen(clientResponse.Token)
		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.FORBIDDEN {
			t.Errorf("expected %q, but got %q", myerrors.FORBIDDEN, myerr.Type)
		}
	})

	t.Run("get all tickets with an invalid token", func(t *testing.T) {
		uc := newTicketUC()

		_, myerr := uc.GetAllOpen("123")
		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.UNAUTHORIZED {
			t.Errorf("expected %q, but got %q", myerrors.UNAUTHORIZED, myerr.Type)
		}
	})
}
