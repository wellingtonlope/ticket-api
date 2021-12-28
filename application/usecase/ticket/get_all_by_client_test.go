package ticket

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
	"testing"
)

func TestGetAllByClient(t *testing.T) {
	t.Run("get all tickets", func(t *testing.T) {
		uc := newTicketUC()
		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password")
		operator.Profile = domain.PROFILE_OPERATOR
		_, _ = uc.UserUseCase.UserRepository.Insert(operator)
		operatorResponse, _ := uc.UserUseCase.Login("operator@mail.com", "password")
		clientResponse, _ := uc.UserUseCase.Register("client", "client@mail.com", "password")
		_, _ = uc.Open("title", "description", operatorResponse.Token)
		_, _ = uc.Open("title", "description", clientResponse.Token)

		got, _ := uc.GetAllByClient(operator.ID, operatorResponse.Token)
		if got == nil {
			t.Error("expected a ticket slice, but got a nil")
		}
		if len(*got) != 1 {
			t.Errorf("expected %q, but got %q", 1, len(*got))
		}
	})

	t.Run("get all tickets without tickets in database", func(t *testing.T) {
		uc := newTicketUC()
		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password")
		operator.Profile = domain.PROFILE_OPERATOR
		_, _ = uc.UserUseCase.UserRepository.Insert(operator)
		operatorResponse, _ := uc.UserUseCase.Login("operator@mail.com", "password")

		got, _ := uc.GetAllByClient(operator.ID, operatorResponse.Token)
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
		client, _ := domain.UserRegister("client", "client@mail.com", "password")
		_, _ = uc.UserUseCase.UserRepository.Insert(client)
		clientResponse, _ := uc.UserUseCase.Login("client@mail.com", "password")

		_, _ = uc.Open("title", "description", operatorResponse.Token)
		_, _ = uc.Open("title", "description", clientResponse.Token)

		got, _ := uc.GetAllByClient(client.ID, clientResponse.Token)
		if got == nil {
			t.Error("expected a ticket slice, but got a nil")
		}
		if len(*got) != 1 {
			t.Errorf("expected %q, but got %q", 1, len(*got))
		}
	})

	t.Run("get all tickets with a client from another user", func(t *testing.T) {
		uc := newTicketUC()
		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password")
		operator.Profile = domain.PROFILE_OPERATOR
		_, _ = uc.UserUseCase.UserRepository.Insert(operator)
		operatorResponse, _ := uc.UserUseCase.Login("operator@mail.com", "password")
		client, _ := domain.UserRegister("client", "client@mail.com", "password")
		_, _ = uc.UserUseCase.UserRepository.Insert(client)
		clientResponse, _ := uc.UserUseCase.Login("client@mail.com", "password")

		_, _ = uc.Open("title", "description", operatorResponse.Token)
		_, _ = uc.Open("title", "description", clientResponse.Token)

		_, myerr := uc.GetAllByClient(operator.ID, clientResponse.Token)
		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.FORBIDDEN {
			t.Errorf("expected %q, but got %q", myerrors.FORBIDDEN, myerr.Type)
		}
	})

	t.Run("get all tickets with an invalid token", func(t *testing.T) {
		uc := newTicketUC()

		_, myerr := uc.GetAll("123")
		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.UNAUTHORIZED {
			t.Errorf("expected %q, but got %q", myerrors.UNAUTHORIZED, myerr.Type)
		}
	})
}
