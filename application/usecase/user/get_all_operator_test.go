package user

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"testing"

	"github.com/wellingtonlope/ticket-api/domain"
)

func TestGetAllOperator(t *testing.T) {
	t.Run("with some registered operators", func(t *testing.T) {
		uc := newUserUC()
		_, _ = uc.Register("client", "client@mail.com", "password")
		operator, _ := domain.UserRegister("operator1", "operator1@mail.com", "password")
		operator.Profile = domain.PROFILE_OPERATOR
		_, _ = uc.UserRepository.Insert(operator)
		token, _ := uc.TokenUseCase.Generate(operator)

		operators, _ := uc.GetAllOperator(token)

		if len(*operators) != 1 {
			t.Errorf("expected %q, but got %q", 1, len(*operators))
		}
	})

	t.Run("without a client", func(t *testing.T) {
		uc := newUserUC()
		clientReponse, _ := uc.Register("client", "client@mail.com", "password")

		_, myerr := uc.GetAllOperator(clientReponse.Token)

		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.FORBIDDEN {
			t.Errorf("expected %q, but got %q", myerrors.FORBIDDEN, myerr.Type)
		}
	})

	t.Run("with an invalid token", func(t *testing.T) {
		uc := newUserUC()

		_, myerr := uc.GetAllOperator("123")

		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.UNAUTHORIZED {
			t.Errorf("expected %q, but got %q", myerrors.UNAUTHORIZED, myerr.Type)
		}
	})
}
