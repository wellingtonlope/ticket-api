package http

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/app/usecase/user"
	"github.com/wellingtonlope/ticket-api/internal/domain"
	httpGo "net/http"
	"testing"
	"time"
)

type mockUCRegister struct {
	mock.Mock
}

func (m *mockUCRegister) Handle(input user.RegisterInput) (*user.Output, error) {
	args := m.Called(input)
	var result *user.Output
	if args.Get(0) != nil {
		result = args.Get(0).(*user.Output)
	}
	return result, args.Error(1)
}

type mockUCLogin struct {
	mock.Mock
}

func (m *mockUCLogin) Handle(input user.LoginInput) (*user.Output, error) {
	args := m.Called(input)
	var result *user.Output
	if args.Get(0) != nil {
		result = args.Get(0).(*user.Output)
	}
	return result, args.Error(1)
}

type mockUCGetAllOperators struct {
	mock.Mock
}

func (m *mockUCGetAllOperators) Handle(input user.GetAllOperatorsInput) (*[]user.Output, error) {
	args := m.Called(input)
	var result *[]user.Output
	if args.Get(0) != nil {
		result = args.Get(0).(*[]user.Output)
	}
	return result, args.Error(1)
}

func TestUserResponseFromUserOutput(t *testing.T) {
	t.Run("should return user response from user output", func(t *testing.T) {
		dateString := "2019-01-01T00:00:00"
		date, _ := time.Parse("2006-01-02T15:04:05", dateString)
		output := user.Output{
			ID:        "id",
			Name:      "name",
			Email:     "mail@mail.com",
			Profile:   "CLIENT",
			CreatedAt: &date,
			UpdatedAt: &date,
		}
		response := userResponseFromUserOutput(output)
		assert.Equal(t, output.ID, response.ID)
		assert.Equal(t, output.Name, response.Name)
		assert.Equal(t, output.Email, response.Email)
		assert.Equal(t, dateString, response.CreatedAt)
		assert.Equal(t, dateString, response.UpdatedAt)
	})
	t.Run("should return an empty user response", func(t *testing.T) {
		response := userResponseFromUserOutput(user.Output{})
		assert.Equal(t, "", response.ID)
		assert.Equal(t, "", response.Name)
		assert.Equal(t, "", response.Email)
		assert.Equal(t, "", response.CreatedAt)
		assert.Equal(t, "", response.UpdatedAt)
	})
}

func TestUsersResponseFromUsersOutput(t *testing.T) {
	t.Run("should return users response from users output", func(t *testing.T) {
		response := usersResponseFromUsersOutput([]user.Output{{}, {}})
		assert.Len(t, *response, 2)
	})
}

func TestUserController_Register(t *testing.T) {
	today := time.Now()
	name, email, password, token := "test", "test@mail.com", "password", "the_token"
	output := &user.Output{
		ID:        "id",
		Name:      name,
		Email:     email,
		Profile:   "CLIENT",
		CreatedAt: &today,
	}

	t.Run("should register a user", func(t *testing.T) {
		request := Request{
			Body: fmt.Sprintf("{\"name\":\"%s\",\"email\":\"%s\",\"password\":\"%s\"}", name, email, password),
		}

		mockAuth := new(mockAuth)
		mockAuth.On("Generate", mock.Anything).Return(token, nil)
		mockUCRegister := new(mockUCRegister)
		mockUCRegister.On("Handle", mock.Anything).Return(output, nil)

		uc := UserController{Authenticator: mockAuth, UCRegister: mockUCRegister}

		response := uc.Register(request)

		assert.Equal(t, httpGo.StatusCreated, response.HttpCode)
		assert.Equal(t, fmt.Sprintf("{\"token\":\"%s\"}", token), response.Body)
		mockUCRegister.AssertCalled(t, "Handle", mock.MatchedBy(func(i interface{}) bool {
			input := i.(user.RegisterInput)
			return name == input.Name && email == input.Email && password == input.Password
		}))
	})

	t.Run("shouldn't register a user with invalid JSON", func(t *testing.T) {
		request := Request{
			Body: fmt.Sprintf("{name\":\"%s\",\"email\":\"%s\",\"password\":\"%s\"}", name, email, password),
		}

		mockAuth := new(mockAuth)
		mockAuth.On("Generate", mock.Anything).Return(token, nil)
		mockUCRegister := new(mockUCRegister)
		mockUCRegister.On("Handle", mock.Anything).Return(output, nil)

		uc := UserController{Authenticator: mockAuth, UCRegister: mockUCRegister}

		response := uc.Register(request)

		assert.Equal(t, httpGo.StatusBadRequest, response.HttpCode)
		assert.Equal(t, "{\"message\":\"invalid json body\"}", response.Body)
	})

	t.Run("shouldn't register a user when use case return an error", func(t *testing.T) {
		request := Request{
			Body: fmt.Sprintf("{\"name\":\"%s\",\"email\":\"invalid_mail\",\"password\":\"%s\"}", name, password),
		}

		mockAuth := new(mockAuth)
		mockAuth.On("Generate", mock.Anything).Return(token, nil)
		mockUCRegister := new(mockUCRegister)
		mockUCRegister.On("Handle", mock.Anything).Return(nil, domain.ErrEmailIsInvalid)

		uc := UserController{Authenticator: mockAuth, UCRegister: mockUCRegister}

		response := uc.Register(request)

		assert.Equal(t, httpGo.StatusBadRequest, response.HttpCode)
		assert.Equal(t, "{\"message\":\"email is invalid\"}", response.Body)
	})

	t.Run("shouldn't register a user when authenticator return an error", func(t *testing.T) {
		request := Request{
			Body: fmt.Sprintf("{\"name\":\"%s\",\"email\":\"invalid_mail\",\"password\":\"%s\"}", name, password),
		}

		mockAuth := new(mockAuth)
		mockAuth.On("Generate", mock.Anything).Return("", errors.New("an error"))
		mockUCRegister := new(mockUCRegister)
		mockUCRegister.On("Handle", mock.Anything).Return(output, nil)

		uc := UserController{Authenticator: mockAuth, UCRegister: mockUCRegister}

		response := uc.Register(request)

		assert.Equal(t, httpGo.StatusInternalServerError, response.HttpCode)
		assert.Equal(t, "{\"message\":\"an error\"}", response.Body)
	})
}

func TestUserController_Login(t *testing.T) {
	today := time.Now()
	name, email, password, token := "test", "test@mail.com", "password", "the_token"
	output := &user.Output{
		ID:        "id",
		Name:      name,
		Email:     email,
		Profile:   "CLIENT",
		CreatedAt: &today,
	}

	t.Run("should login a user", func(t *testing.T) {
		request := Request{
			Body: fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\"}", email, password),
		}

		mockAuth := new(mockAuth)
		mockAuth.On("Generate", mock.Anything).Return(token, nil)
		mockUCLogin := new(mockUCLogin)
		mockUCLogin.On("Handle", mock.Anything).Return(output, nil)

		uc := UserController{Authenticator: mockAuth, UCLogin: mockUCLogin}

		response := uc.Login(request)

		assert.Equal(t, httpGo.StatusOK, response.HttpCode)
		assert.Equal(t, fmt.Sprintf("{\"token\":\"%s\"}", token), response.Body)
		mockUCLogin.AssertCalled(t, "Handle", mock.MatchedBy(func(i interface{}) bool {
			input := i.(user.LoginInput)
			return email == input.Email && password == input.Password
		}))
	})

	t.Run("shouldn't login a user with invalid JSON", func(t *testing.T) {
		request := Request{
			Body: fmt.Sprintf("{email\":\"%s\",\"password\":\"%s\"}", email, password),
		}

		mockAuth := new(mockAuth)
		mockAuth.On("Generate", mock.Anything).Return(token, nil)
		mockUCLogin := new(mockUCLogin)
		mockUCLogin.On("Handle", mock.Anything).Return(output, nil)

		uc := UserController{Authenticator: mockAuth, UCLogin: mockUCLogin}

		response := uc.Login(request)

		assert.Equal(t, httpGo.StatusBadRequest, response.HttpCode)
		assert.Equal(t, "{\"message\":\"invalid json body\"}", response.Body)
	})

	t.Run("shouldn't login a user when use case return an error", func(t *testing.T) {
		request := Request{
			Body: fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\"}", email, password),
		}

		mockAuth := new(mockAuth)
		mockAuth.On("Generate", mock.Anything).Return(token, nil)
		mockUCLogin := new(mockUCLogin)
		mockUCLogin.On("Handle", mock.Anything).Return(nil, user.ErrUserEmailPasswordWrong)

		uc := UserController{Authenticator: mockAuth, UCLogin: mockUCLogin}

		response := uc.Login(request)

		assert.Equal(t, httpGo.StatusBadRequest, response.HttpCode)
		assert.Equal(t, "{\"message\":\"email or password is wrong\"}", response.Body)
	})

	t.Run("shouldn't login a user when authenticator return an error", func(t *testing.T) {
		request := Request{
			Body: fmt.Sprintf("{\"name\":\"%s\",\"email\":\"invalid_mail\",\"password\":\"%s\"}", name, password),
		}

		mockAuth := new(mockAuth)
		mockAuth.On("Generate", mock.Anything).Return("", errors.New("an error"))
		mockUCLogin := new(mockUCLogin)
		mockUCLogin.On("Handle", mock.Anything).Return(output, nil)

		uc := UserController{Authenticator: mockAuth, UCLogin: mockUCLogin}

		response := uc.Login(request)

		assert.Equal(t, httpGo.StatusInternalServerError, response.HttpCode)
		assert.Equal(t, "{\"message\":\"an error\"}", response.Body)
	})
}

func TestUserController_GetAllOperators(t *testing.T) {
	date, _ := time.Parse("2006-01-02T15:04:05", "2019-01-01T00:00:00")
	loggedUser := &security.User{
		ID:      "id",
		Name:    "test",
		Profile: "OPERATOR",
	}
	outputs := &[]user.Output{{
		ID:        "id",
		Name:      "test",
		Email:     "test@mail.com",
		Profile:   "CLIENT",
		CreatedAt: &date,
		UpdatedAt: &date,
	}}

	t.Run("should get all operators", func(t *testing.T) {
		request := Request{
			LoggedUser: loggedUser,
		}

		mockUCGetAllOperators := new(mockUCGetAllOperators)
		mockUCGetAllOperators.On("Handle", mock.Anything).Return(outputs, nil)

		uc := UserController{UCGetAllOperators: mockUCGetAllOperators}

		response := uc.GetAllOperators(request)

		assert.Equal(t, httpGo.StatusOK, response.HttpCode)
		assert.Equal(t,
			"[{\"id\":\"id\",\"name\":\"test\",\"email\":\"test@mail.com\","+
				"\"created_at\":\"2019-01-01T00:00:00\",\"updated_at\":\"2019-01-01T00:00:00\"}]",
			response.Body)
	})

	t.Run("shouldn't get all operators with a client", func(t *testing.T) {
		request := Request{
			LoggedUser: loggedUser,
		}

		mockUCGetAllOperators := new(mockUCGetAllOperators)
		mockUCGetAllOperators.On("Handle", mock.Anything).Return(nil, security.ErrForbidden)

		uc := UserController{UCGetAllOperators: mockUCGetAllOperators}

		response := uc.GetAllOperators(request)

		assert.Equal(t, httpGo.StatusForbidden, response.HttpCode)
		assert.Equal(t, "{\"message\":\"user don't have permission\"}", response.Body)
	})
}
