package http

import (
	"encoding/json"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/app/usecase/user"
	"github.com/wellingtonlope/ticket-api/internal/domain"
	"net/http"
	"time"
)

type (
	UserRegisterRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	UserLoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	UserTokenResponse struct {
		Token string `json:"token"`
	}
	UserResponse struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at,omitempty"`
	}
)

func userResponseFromUserOutput(user *user.UserOutput) *UserResponse {
	var updatedAt string
	if user.UpdatedAt != nil {
		updatedAt = user.UpdatedAt.Format("2006-01-02 15:04:05")
	}
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: updatedAt,
	}
}

func usersResponseFromUsersOutput(outputs *[]user.UserOutput) *[]UserResponse {
	response := make([]UserResponse, 0, len(*outputs))
	for _, output := range *outputs {
		response = append(response, *userResponseFromUserOutput(&output))
	}
	return &response
}

type UserController struct {
	UCRegister        *user.Register
	UCLogin           *user.Login
	UCGetAllOperators *user.GetAllOperators
	Authenticator     security.Authenticator
}

func (c *UserController) Register(request Request) Response {
	var registerRequest UserRegisterRequest
	err := json.Unmarshal([]byte(request.Body), &registerRequest)
	if err != nil {
		return Response{
			HttpCode: http.StatusBadRequest,
			Body:     wrapError(ErrInvalidJsonBody),
		}
	}

	input := user.RegisterInput{
		Name:      registerRequest.Name,
		Email:     registerRequest.Email,
		Password:  registerRequest.Password,
		CreatedAt: time.Now(),
	}
	output, err := c.UCRegister.Handle(input)
	if err != nil {
		httpStatus := http.StatusInternalServerError
		switch err {
		case user.ErrUserAlreadyExists:
		case domain.ErrNameIsInvalid:
		case domain.ErrEmailIsInvalid:
		case domain.ErrPasswordIsInvalid:
			httpStatus = http.StatusBadRequest
		}
		return Response{
			HttpCode: httpStatus,
			Body:     wrapError(err),
		}
	}

	token, err := c.Authenticator.Generate(security.User{
		ID:      output.ID,
		Name:    output.Name,
		Profile: output.Profile,
	})
	if err != nil {
		return Response{
			HttpCode: http.StatusInternalServerError,
			Body:     wrapError(err),
		}
	}

	return Response{
		HttpCode: http.StatusCreated,
		Body:     wrapBody(UserTokenResponse{Token: token}),
	}
}

func (c *UserController) Login(request Request) Response {
	var loginRequest UserLoginRequest
	err := json.Unmarshal([]byte(request.Body), &loginRequest)
	if err != nil {
		return Response{
			HttpCode: http.StatusBadRequest,
			Body:     wrapError(ErrInvalidJsonBody),
		}
	}

	input := user.LoginInput{
		Email:    loginRequest.Email,
		Password: loginRequest.Password,
	}
	output, err := c.UCLogin.Handle(input)
	if err != nil {
		httpStatus := http.StatusInternalServerError
		switch err {
		case user.ErrUserEmailPasswordWrong:
			httpStatus = http.StatusBadRequest
		}
		return Response{
			HttpCode: httpStatus,
			Body:     wrapError(err),
		}
	}

	token, err := c.Authenticator.Generate(security.User{
		ID:      output.ID,
		Name:    output.Name,
		Profile: output.Profile,
	})
	if err != nil {
		return Response{
			HttpCode: http.StatusInternalServerError,
			Body:     wrapError(err),
		}
	}

	return Response{
		HttpCode: http.StatusOK,
		Body:     wrapBody(UserTokenResponse{Token: token}),
	}
}

func (c *UserController) GetAllOperators(request Request) Response {
	output, err := c.UCGetAllOperators.Handle(user.GetAllOperatorsInput{
		LoggedUser: *request.LoggedUser,
	})
	if err != nil {
		httpStatus := http.StatusInternalServerError
		switch err {
		case security.ErrForbidden:
			httpStatus = http.StatusForbidden
		}
		return Response{
			HttpCode: httpStatus,
			Body:     wrapError(err),
		}
	}

	return Response{
		HttpCode: http.StatusCreated,
		Body:     wrapBody(usersResponseFromUsersOutput(output)),
	}
}